import md5 from 'js-md5'

/**
 * Legacy Angular user modal sends `newPassword` as MD5 hex (uppercase).
 * {@link com.hmdm.util.PasswordUtil#getHashFromMd5} then derives stored hash from that value.
 */
export function encodePasswordForUserSave(plainPassword: string): string {
  return (md5 as unknown as (s: string) => string)(plainPassword).toUpperCase()
}
