import { HttpContextContract } from '@ioc:Adonis/Core/HttpContext'

export default class Deobfuscator {
  public async handle ({ request }: HttpContextContract, next: () => Promise<void>) {
    let obfuscatedBody = ''
    for (const key of Object.keys(request.body()))
      obfuscatedBody += Buffer.from(request.body()[key], 'base64').toString('utf-8')

    const newBody = (() => {
      try {
        return JSON.parse(obfuscatedBody)
      }
      catch {
        return null
      }
    })()

    if (newBody)
      request.updateBody(newBody)

    await next()
  }
}
