import JSEncrypt from 'jsencrypt'
import md5 from 'js-md5'

export const md5Hash = md5 as unknown as (value: string) => string

export function md5UpperHex(value: string): string {
  return md5Hash(value).toUpperCase()
}

function pkixBase64ToPem(base64: string): string {
  const lines = base64.match(/.{1,64}/g) ?? [base64]
  return `-----BEGIN PUBLIC KEY-----\n${lines.join('\n')}\n-----END PUBLIC KEY-----`
}

/**
 * Same login encoding as legacy Angular app:
 * - if publicKey exists: RSA encrypt password
 * - else: MD5 uppercase
 */
export function encodeLoginPassword(
  rawPassword: string,
  publicKeyPkixBase64: string | null | undefined
): string {
  if (publicKeyPkixBase64 && publicKeyPkixBase64.trim()) {
    const encryptor = new JSEncrypt()
    encryptor.setPublicKey(pkixBase64ToPem(publicKeyPkixBase64.trim()))
    const encrypted = encryptor.encrypt(rawPassword)
    if (!encrypted) {
      throw new Error('Password encryption failed.')
    }
    return encrypted
  }
  return md5UpperHex(rawPassword)
}
