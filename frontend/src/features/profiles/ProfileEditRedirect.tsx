import { Navigate, useParams, useSearchParams } from 'react-router-dom'

/** Legacy routes → workspace URL (020). */
export function ProfileEditRedirect() {
  const { profileId, versionId } = useParams<{ profileId: string; versionId?: string }>()
  const [searchParams] = useSearchParams()
  const id = profileId ?? ''
  if (!id) {
    return <Navigate to="/profiles" replace />
  }
  const qs = new URLSearchParams(searchParams)
  qs.set('open', id)
  qs.set('section', 'editor')
  if (versionId) {
    qs.set('versionId', versionId)
  } else {
    qs.delete('versionId')
  }
  qs.delete('readOnly')
  return <Navigate to={`/profiles?${qs.toString()}`} replace />
}
