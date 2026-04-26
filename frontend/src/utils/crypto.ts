import bcrypt from 'bcryptjs'

/**
 * 密码加密
 * @param password 原始密码
 * @returns 加密后的密码
 */
export const encryptPassword = (password: string): string => {
  // 使用bcrypt加密，生成盐值并加密密码
  const salt = bcrypt.genSaltSync(10)
  const hash = bcrypt.hashSync(password, salt)
  return hash
}
