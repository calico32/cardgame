import type { NextPage } from 'next'
import Head from 'next/head'
import { Fragment, useCallback, useState } from 'react'
import Wrapper from '../components/Wrapper'

const Home: NextPage = () => {
  const [requestMethod, setRequestMethod] = useState('GET')
  const [requestURL, setRequestURL] = useState('http://localhost:8080/api')
  const [password, setPassword] = useState('')
  const [response, setResponse] = useState('')
  const [ws, setWs] = useState<WebSocket | null>(null)
  const [payload, setPayload] = useState<string>('')
  const [messages, setMessages] = useState<string[]>([])

  const request = useCallback(async () => {
    try {
      const res = await fetch(requestURL, {
        method: requestMethod,
        headers: {
          'X-Password': password,
        },
      })
      setResponse(JSON.stringify(await res.json()))
    } catch (err: any) {
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      setResponse(err.stack)
    }
  }, [requestURL, requestMethod, password])

  const connectWs = useCallback(() => {
    const ws = new WebSocket(requestURL.replace(/^http/, 'ws') + '?password=' + password)
    ws.onopen = () => {
      console.log('connected')
      setWs(ws)
    }
    ws.onmessage = (e) => {
      setMessages((messages) => [...messages, e.data])
    }
  }, [requestURL, password])

  const send = useCallback(() => {
    if (ws) {
      try {
        const json = JSON.parse(payload)
        ws.send(json)
      } catch (e) {
        console.error(e)
      }
    } else {
      console.error('ws not connected')
    }
  }, [ws, payload])

  return (
    <>
      <Head>
        <title>App Name</title>
      </Head>
      <Wrapper page>
        <select value={requestMethod} onChange={(e) => setRequestMethod(e.target.value)}>
          <option value="GET">GET</option>
          <option value="POST">POST</option>
        </select>
        <input value={requestURL} onChange={(e) => setRequestURL(e.target.value)} />
        <input value={password} onChange={(e) => setPassword(e.target.value)} />
        <button onClick={request}>Send Request</button>
        <button onClick={connectWs}>Connect WebSocket</button>
        {response && <pre>{response}</pre>}
        <br />
        <br />
        <br />
        <br />
        <textarea value={payload} onChange={(e) => setPayload(e.target.value)} />
        <br />
        <button onClick={send}>Send</button>
        {messages.map((message, i) => (
          <Fragment key={i}>
            <pre>{message}</pre>
            <br />
            <br />
          </Fragment>
        ))}
      </Wrapper>
    </>
  )
}

export default Home
