/** Mirrors `password.service.js` `checkQuality`. */
export function checkPasswordQuality(password: string, length: number, strength: number): boolean {
  if (password.length < length) return false
  if (strength === 0) return true
  const hasDigits = /\d/.test(password)
  const hasLower = /[a-z]/.test(password)
  const hasCaps = /[A-Z]/.test(password)
  if (strength === 1) {
    return hasDigits && hasLower && hasCaps
  }
  if (strength === 2) {
    const hasSpecial = /[_\-.,!#$%()=+;*/]/.test(password)
    return hasDigits && hasLower && hasCaps && hasSpecial
  }
  return false
}

export function passwordPolicyHint(length: number, strength: number): string {
  const bits: string[] = []
  if (length > 0) bits.push(`At least ${length} characters`)
  if (strength === 1) bits.push('digits, upper and lower case letters')
  if (strength === 2) bits.push('digits, mixed case, and special chars from _-. ,!#$%()=+;*/ ')
  return bits.join('; ') || 'No extra rules.'
}
