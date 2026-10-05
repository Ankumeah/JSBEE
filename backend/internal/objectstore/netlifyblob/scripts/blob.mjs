import { getStore } from '@netlify/blobs'

const EXIT_NOT_FOUND = 3

function fail(code, message, exit = 1) {
  process.stderr.write(JSON.stringify({ code, message }) + '\n')
  process.exit(exit)
}

function checkKey(key) {
  if (typeof key !== 'string' || key.length === 0 || key.length > 220) {
    fail('bad_key', 'invalid key')
  }
  if (!/^[A-Za-z0-9._/-]+$/.test(key) || key.includes('..')) {
    fail('bad_key', 'invalid key')
  }
}

async function readStdin() {
  const chunks = []
  for await (const chunk of process.stdin) {
    chunks.push(chunk)
  }
  return Buffer.concat(chunks)
}

function store() {
  const siteID = process.env.BLOB_SITE_ID
  const token = process.env.BLOB_TOKEN
  const apiURL = process.env.BLOB_API_URL
  const name = process.env.BLOB_STORE
  if (!siteID || !token || !name) {
    fail('config', 'missing blob credentials')
  }

  return getStore({ name, siteID, token, apiURL, consistency: 'strong' })
}

async function main() {
  const [op, a, b] = process.argv.slice(2)
  const s = store()

  if (op === 'ping') {
    await s.list({ prefix: '__jsbee_ping__', paginate: false })
    process.stdout.write('ok\n')
    return
  }

  if (op === 'set') {
    checkKey(a)
    const body = await readStdin()
    await s.set(a, body)
    process.stdout.write(JSON.stringify({ size: body.length }) + '\n')
    return
  }

  if (op === 'get') {
    checkKey(a)
    const data = await s.get(a, { type: 'arrayBuffer' })
    if (data === null) {
      fail('not_found', 'no such blob', EXIT_NOT_FOUND)
    }
    process.stdout.write(Buffer.from(data))
    return
  }

  if (op === 'meta') {
    checkKey(a)
    const entry = await s.getWithMetadata(a, { type: 'arrayBuffer' })
    if (entry === null || entry.data === null) {
      fail('not_found', 'no such blob', EXIT_NOT_FOUND)
    }
    process.stdout.write(
      JSON.stringify({ size: entry.data.byteLength, etag: entry.etag }) + '\n',
    )
    return
  }

  if (op === 'copy' || op === 'move') {
    checkKey(a)
    checkKey(b)

    const data = await s.get(a, { type: 'arrayBuffer' })
    if (data === null) {
      fail('not_found', 'no such blob', EXIT_NOT_FOUND)
    }
    await s.set(b, Buffer.from(data))
    if (op === 'move') {
      await s.delete(a)
    }
    process.stdout.write('ok\n')
    return
  }

  if (op === 'del') {
    checkKey(a)

    await s.delete(a)
    process.stdout.write('ok\n')
    return
  }

  if (op === 'list') {
    const prefix = typeof a === 'string' ? a : ''
    const { blobs } = await s.list({ prefix })
    process.stdout.write(
      JSON.stringify({ keys: blobs.map((e) => e.key) }) + '\n',
    )
    return
  }

  fail('bad_op', `unknown op ${op}`)
}

main().catch((err) => {
  fail('store', err && err.message ? err.message : String(err))
})
