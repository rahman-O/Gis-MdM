package domain

import "strings"

// Transliterate converts customer name to org-admin login (Java CustomerDAO.transliterate).
func Transliterate(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, c := range s {
		switch c {
		case 'а':
			b.WriteByte('a')
		case 'б':
			b.WriteByte('b')
		case 'в':
			b.WriteByte('v')
		case 'г':
			b.WriteByte('g')
		case 'д':
			b.WriteByte('d')
		case 'е', 'ё':
			b.WriteByte('e')
		case 'ж':
			b.WriteString("zh")
		case 'з':
			b.WriteByte('z')
		case 'и', 'й':
			b.WriteByte('i')
		case 'к':
			b.WriteByte('k')
		case 'л':
			b.WriteByte('l')
		case 'м':
			b.WriteByte('m')
		case 'н':
			b.WriteByte('n')
		case 'о':
			b.WriteByte('o')
		case 'п':
			b.WriteByte('p')
		case 'р':
			b.WriteByte('r')
		case 'с':
			b.WriteByte('s')
		case 'т':
			b.WriteByte('t')
		case 'у':
			b.WriteByte('u')
		case 'ф':
			b.WriteByte('f')
		case 'х':
			b.WriteString("kh")
		case 'ц':
			b.WriteString("ts")
		case 'ч':
			b.WriteString("ch")
		case 'ш':
			b.WriteString("sh")
		case 'щ':
			b.WriteString("shch")
		case 'ъ':
			b.WriteString("ie")
		case 'ы':
			b.WriteByte('y')
		case 'ь':
			b.WriteByte('-')
		case 'э':
			b.WriteByte('e')
		case 'ю':
			b.WriteString("yu")
		case 'я':
			b.WriteString("ya")
		case ' ':
			b.WriteByte('_')
		default:
			b.WriteRune(c)
		}
	}
	return b.String()
}
