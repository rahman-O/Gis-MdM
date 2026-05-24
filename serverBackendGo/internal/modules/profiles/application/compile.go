package application

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	cfgdomain "github.com/gis-mdm/server-backend-go/internal/modules/configurations/domain"
	"github.com/gis-mdm/server-backend-go/internal/modules/profiles/domain"
	syncapp "github.com/gis-mdm/server-backend-go/internal/modules/sync/application"
	syncdomain "github.com/gis-mdm/server-backend-go/internal/modules/sync/domain"
	"github.com/gis-mdm/server-backend-go/internal/platform/storage"
)

// CompileInput loads a profile version into a publish-time artifact.
type CompileInput struct {
	ProfileID        int
	ProfileVersionID int
	VersionNumber    int
	Payload          cfgdomain.Configuration
	SettingsJSON     []byte
}

// ArtifactCompiler builds profile_version_artifacts from version rows.
type ArtifactCompiler struct {
	db         *sql.DB
	baseURL    string
	filesDir   string
}

func NewArtifactCompiler(db *sql.DB, baseURL, filesDir string) *ArtifactCompiler {
	return &ArtifactCompiler{db: db, baseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"), filesDir: filesDir}
}

func (c *ArtifactCompiler) Compile(ctx context.Context, in CompileInput) (*domain.ProfileArtifact, string, error) {
	apps, err := c.loadApplications(ctx, in.ProfileVersionID)
	if err != nil {
		return nil, "", err
	}
	files, err := c.loadFiles(ctx, in.ProfileVersionID)
	if err != nil {
		return nil, "", err
	}
	cfgSettings, err := c.loadConfigAppSettings(ctx, in.ProfileVersionID)
	if err != nil {
		return nil, "", err
	}
	permissive := false
	if in.Payload.Permissive != nil {
		permissive = *in.Payload.Permissive
	}
	password := ""
	if in.Payload.Password != nil {
		password = *in.Payload.Password
	}
	artifact := &domain.ProfileArtifact{
		ProfileID:                 in.ProfileID,
		ProfileVersionID:          in.ProfileVersionID,
		VersionNumber:             in.VersionNumber,
		Password:                  password,
		BackgroundColor:           in.Payload.BackgroundColor,
		TextColor:                 in.Payload.TextColor,
		BackgroundImageURL:        in.Payload.BackgroundImageURL,
		Permissive:                permissive,
		SettingsJSON:              in.SettingsJSON,
		Applications:              apps,
		Files:                     files,
		ConfigApplicationSettings: cfgSettings,
	}
	raw, err := json.Marshal(artifact)
	if err != nil {
		return nil, "", err
	}
	sum := sha256.Sum256(raw)
	return artifact, hex.EncodeToString(sum[:]), nil
}

func (c *ArtifactCompiler) loadApplications(ctx context.Context, versionID int) ([]syncdomain.SyncApplication, error) {
	rows, err := c.db.QueryContext(ctx, `
		SELECT a.id, a.name, a.pkg, COALESCE(av.version, ''), COALESCE(av.url, ''),
			COALESCE(a.type, 'app'), COALESCE(av.urlarm64, ''), COALESCE(av.urlarmeabi, ''), COALESCE(av.split, false),
			av.versioncode,
			a.runafterinstall, a.runatboot, COALESCE(a.system, false), COALESCE(a.usekiosk, false),
			a.icontext, a.intent,
			COALESCE(pva.showicon, a.showicon, true), COALESCE(pva.remove, false),
			pva.screenorder, pva.keycode, COALESCE(pva.bottom, false), COALESCE(pva.longtap, false),
			COALESCE(pap.value = 'true', false),
			uf.filepath
		FROM profile_version_applications pva
		JOIN applications a ON a.id = pva.applicationid
		LEFT JOIN applicationversions av ON av.id = pva.applicationversionid
		LEFT JOIN profile_version_application_parameters pap
			ON pap.profile_version_id = pva.profile_version_id
		   AND pap.applicationid = pva.applicationid AND pap.name = 'skipVersionCheck'
		LEFT JOIN icons ON icons.id = a.iconid
		LEFT JOIN uploadedfiles uf ON uf.id = icons.fileid
		WHERE pva.profile_version_id = $1
		ORDER BY pva.screenorder NULLS LAST, a.name`, versionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]syncdomain.SyncApplication, 0)
	screenIdx := 0
	for rows.Next() {
		var app syncdomain.SyncApplication
		var urlArm64, urlArm, iconPath, iconText, intent sql.NullString
		var split, runAfter, runAtBoot, system, useKiosk, showIcon, remove, bottom, longTap, skipVer bool
		var verCode sql.NullInt64
		var screenOrder, keyCode sql.NullInt64
		if err := rows.Scan(
			&app.ID, &app.Name, &app.Pkg, &app.Version, &app.URL, &app.Type,
			&urlArm64, &urlArm, &split,
			&verCode, &runAfter, &runAtBoot, &system, &useKiosk,
			&iconText, &intent,
			&showIcon, &remove, &screenOrder, &keyCode, &bottom, &longTap, &skipVer,
			&iconPath,
		); err != nil {
			return nil, err
		}
		if split {
			if urlArm64.Valid && urlArm64.String != "" {
				app.URL = urlArm64.String
			} else if urlArm.Valid && urlArm.String != "" {
				app.URL = urlArm.String
			}
		}
		if verCode.Valid && verCode.Int64 > 0 {
			code := int(verCode.Int64)
			app.Code = &code
		}
		if iconPath.Valid && strings.TrimSpace(iconPath.String) != "" {
			u := storage.BuildPublicURL(c.baseURL, c.filesDir, iconPath.String)
			app.Icon = &u
		}
		app.ShowIcon = boolTrueOnly(showIcon)
		if screenOrder.Valid {
			s := int(screenOrder.Int64)
			app.ScreenOrder = &s
		} else {
			screenIdx++
			s := screenIdx
			app.ScreenOrder = &s
		}
		app.UseKiosk = boolTrueOnly(useKiosk)
		app.Remove = boolTrueOnly(remove)
		app.System = boolTrueOnly(system)
		app.RunAfterInstall = boolTrueOnly(runAfter)
		app.RunAtBoot = boolTrueOnly(runAtBoot)
		app.SkipVersion = boolTrueOnly(skipVer)
		app.Bottom = boolTrueOnly(bottom)
		app.LongTap = boolTrueOnly(longTap)
		if iconText.Valid && strings.TrimSpace(iconText.String) != "" {
			t := iconText.String
			app.IconText = &t
		}
		if keyCode.Valid && keyCode.Int64 > 0 {
			k := int(keyCode.Int64)
			app.KeyCode = &k
		}
		if intent.Valid && strings.TrimSpace(intent.String) != "" {
			in := intent.String
			app.Intent = &in
		}
		out = append(out, app)
	}
	return out, rows.Err()
}

func (c *ArtifactCompiler) loadFiles(ctx context.Context, versionID int) ([]syncdomain.SyncConfigurationFile, error) {
	rows, err := c.db.QueryContext(ctx, `
		SELECT COALESCE(path, ''), COALESCE(externalurl, ''), COALESCE(url, ''), remove
		FROM profile_version_files WHERE profile_version_id = $1`, versionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]syncdomain.SyncConfigurationFile, 0)
	for rows.Next() {
		var path, extURL, url string
		var remove bool
		if err := rows.Scan(&path, &extURL, &url, &remove); err != nil {
			return nil, err
		}
		f := syncdomain.SyncConfigurationFile{Path: path, Remove: remove}
		if extURL != "" {
			f.URL = extURL
			f.External = true
		} else if path != "" {
			f.URL = storage.BuildPublicURL(c.baseURL, c.filesDir, path)
		} else if url != "" {
			f.URL = url
		}
		if f.Path == "" || (!remove && f.URL == "") {
			continue
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func (c *ArtifactCompiler) loadConfigAppSettings(ctx context.Context, versionID int) ([]syncdomain.SyncApplicationSetting, error) {
	rows, err := c.db.QueryContext(ctx, `
		SELECT COALESCE(a.pkg, ''), s.name, s.type, s.value
		FROM profile_version_application_settings s
		LEFT JOIN applications a ON a.id = s.applicationid
		WHERE s.profile_version_id = $1`, versionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]syncdomain.SyncApplicationSetting, 0)
	for rows.Next() {
		var s syncdomain.SyncApplicationSetting
		var typ sql.NullString
		if err := rows.Scan(&s.PackageID, &s.Name, &typ, &s.Value); err != nil {
			return nil, err
		}
		s.Type = parseSettingType(typ)
		out = append(out, s)
	}
	return out, rows.Err()
}

func boolTrueOnly(v bool) *bool {
	if v {
		return &v
	}
	return nil
}

func parseSettingType(typ sql.NullString) int {
	if !typ.Valid {
		return 0
	}
	var n int
	if _, err := fmt.Sscanf(typ.String, "%d", &n); err == nil {
		return n
	}
	return 0
}

// ApplyArtifactToSyncResponse maps a compiled artifact onto a sync response (parity with sync_configuration_mapper).
func ApplyArtifactToSyncResponse(resp *syncdomain.SyncResponse, artifact *domain.ProfileArtifact) {
	if resp == nil || artifact == nil {
		return
	}
	var bgi *string
	if artifact.BackgroundImageURL != nil && *artifact.BackgroundImageURL != "" {
		bgi = artifact.BackgroundImageURL
	}
	syncapp.ApplyConfigurationPolicy(resp, artifact.SettingsJSON, bgi)
	resp.BackgroundColor = artifact.BackgroundColor
	resp.TextColor = artifact.TextColor
	if artifact.BackgroundImageURL != nil {
		resp.BackgroundImageURL = artifact.BackgroundImageURL
	}
	permissive := artifact.Permissive
	resp.Permissive = &permissive
	resp.Applications = artifact.Applications
	resp.Files = artifact.Files
}
