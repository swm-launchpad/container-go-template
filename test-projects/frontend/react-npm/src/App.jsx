import { useState } from 'react'
import './App.css'

function App() {
  const [counter, setCounter] = useState(0)
  const [lastActionTime, setLastActionTime] = useState('No actions yet')
  const [newItem, setNewItem] = useState('')
  const [items, setItems] = useState([])

  const updateTimestamp = () => {
    setLastActionTime(new Date().toLocaleString())
  }

  const increment = () => {
    setCounter(counter + 1)
    updateTimestamp()
  }

  const decrement = () => {
    setCounter(counter - 1)
    updateTimestamp()
  }

  const reset = () => {
    setCounter(0)
    updateTimestamp()
  }

  const addItem = () => {
    if (newItem.trim()) {
      setItems([...items, newItem.trim()])
      setNewItem('')
      updateTimestamp()
    }
  }

  const removeItem = (index) => {
    setItems(items.filter((_, i) => i !== index))
    updateTimestamp()
  }

  const handleKeyPress = (e) => {
    if (e.key === 'Enter') {
      addItem()
    }
  }

  const getStatus = () => {
    if (counter > 0) return { text: 'Positive', className: 'positive' }
    if (counter < 0) return { text: 'Negative', className: 'negative' }
    return { text: 'Ready', className: 'ready' }
  }

  const status = getStatus()

  return (
    <div id="app">
      <h1>React Test Application</h1>

      {/* Counter Section */}
      <div className="section">
        <h2>Counter</h2>
        <div className="counter-display">
          <span className="count">{counter}</span>
        </div>
        <div className="button-group">
          <button onClick={increment}>Increment</button>
          <button onClick={decrement}>Decrement</button>
          <button onClick={reset}>Reset</button>
        </div>
        <div className={`status ${status.className}`}>
          Status: {status.text}
        </div>
      </div>

      {/* Timestamp Section */}
      <div className="section">
        <h2>Last Action</h2>
        <p className="timestamp">{lastActionTime}</p>
      </div>

      {/* Items List Section */}
      <div className="section">
        <h2>Items List</h2>
        <div className="input-group">
          <input
            type="text"
            value={newItem}
            onChange={(e) => setNewItem(e.target.value)}
            onKeyPress={handleKeyPress}
            placeholder="Enter item name"
          />
          <button onClick={addItem}>Add Item</button>
        </div>
        <ul className="items-list">
          {items.map((item, index) => (
            <li key={index}>
              {item}
              <button onClick={() => removeItem(index)} className="remove-btn">
                Remove
              </button>
            </li>
          ))}
        </ul>
        {items.length === 0 && (
          <p className="empty-message">No items yet. Add some!</p>
        )}
      </div>
    </div>
  )
}

export default App
