import React, { useEffect, useState } from 'react'
import './App.css'

function App() {
  const userId = 2

  const [messages, setMessages] = useState([])
  const [userMessage, setUserMessage] = useState('')

  const [firstLoadFlag, setFirstLoadFlag] = useState(true)
  const [conversationId, setConversationId] = useState(null)

  const handleGetConversation = async () => {
    const response = await fetch(
      `http://localhost:8080/conversation/${conversationId}`,
    )

    console.log(response.json())
  }

  useEffect(() => {
    if (firstLoadFlag && conversationId !== null) {
      handleGetConversation()
    }
    setFirstLoadFlag(false)
  }, [conversationId, firstLoadFlag])

  // Function to handle the submission of the user's message
  const sendMessage = async (e) => {
    let currentConversationId = conversationId
    e.preventDefault()
    if (!userMessage.trim()) return

    const newMessage = { author: 'User', text: userMessage }

    if (!conversationId) {
      const newConversationResponse = await fetch(
        'http://localhost:8080/conversation',
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ userId }),
        },
      )

      if (newConversationResponse.ok) {
        const {
          data: { ID: newConversationId },
        } = await newConversationResponse.json()
        currentConversationId = newConversationId
        setConversationId(newConversationId)
      }
    }

    setMessages([...messages, newMessage])

    // Send the user's message to the backend and await the response
    const response = await fetch(
      `http://localhost:8080/conversation/${currentConversationId}/message`,
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          content: userMessage,
          userId,
          conversationId: currentConversationId,
        }),
      },
    )

    if (response.ok) {
      const { reply } = await response.json()
      setMessages((prevMessages) => [
        ...prevMessages,
        { author: 'Backend', text: reply },
      ])
    }

    // Clear the input after sending
    setUserMessage('')
  }

  return (
    <div className="App">
      {/* <header className="App-header">
        <h1>Simple Chat</h1>
      </header> */}
      <div className="chat-view">
        {messages.map((msg, index) => (
          <div key={index} className={`message ${msg.author}`}>
            <span>{msg.author}: </span>
            {msg.text}
          </div>
        ))}
      </div>
      <form className="message-form" onSubmit={sendMessage}>
        <input
          type="text"
          value={userMessage}
          onChange={(e) => setUserMessage(e.target.value)}
          placeholder="Type a message..."
        />
        <button type="submit">Send</button>
      </form>
    </div>
  )
}

export default App
