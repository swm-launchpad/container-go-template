<template>
  <div id="app">
    <h1>Vue.js Test Application</h1>

    <!-- Counter Section -->
    <div class="section">
      <h2>Counter</h2>
      <div class="counter-display">
        <span class="count">{{ counter }}</span>
      </div>
      <div class="button-group">
        <button @click="increment">Increment</button>
        <button @click="decrement">Decrement</button>
        <button @click="reset">Reset</button>
      </div>
      <div class="status" :class="statusClass">
        Status: {{ status }}
      </div>
    </div>

    <!-- Timestamp Section -->
    <div class="section">
      <h2>Last Action</h2>
      <p class="timestamp">{{ lastActionTime }}</p>
    </div>

    <!-- Items List Section -->
    <div class="section">
      <h2>Items List</h2>
      <div class="input-group">
        <input
          v-model="newItem"
          @keyup.enter="addItem"
          placeholder="Enter item name"
          type="text"
        />
        <button @click="addItem">Add Item</button>
      </div>
      <ul class="items-list">
        <li v-for="(item, index) in items" :key="index">
          {{ item }}
          <button @click="removeItem(index)" class="remove-btn">Remove</button>
        </li>
      </ul>
      <p v-if="items.length === 0" class="empty-message">No items yet. Add some!</p>
    </div>
  </div>
</template>

<script>
export default {
  name: 'App',
  data() {
    return {
      counter: 0,
      lastActionTime: 'No actions yet',
      newItem: '',
      items: []
    }
  },
  computed: {
    status() {
      if (this.counter > 0) return 'Positive'
      if (this.counter < 0) return 'Negative'
      return 'Ready'
    },
    statusClass() {
      if (this.counter > 0) return 'positive'
      if (this.counter < 0) return 'negative'
      return 'ready'
    }
  },
  methods: {
    increment() {
      this.counter++
      this.updateTimestamp()
    },
    decrement() {
      this.counter--
      this.updateTimestamp()
    },
    reset() {
      this.counter = 0
      this.updateTimestamp()
    },
    addItem() {
      if (this.newItem.trim()) {
        this.items.push(this.newItem.trim())
        this.newItem = ''
        this.updateTimestamp()
      }
    },
    removeItem(index) {
      this.items.splice(index, 1)
      this.updateTimestamp()
    },
    updateTimestamp() {
      this.lastActionTime = new Date().toLocaleString()
    }
  }
}
</script>

<style scoped>
#app {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  font-family: Arial, sans-serif;
}

h1 {
  color: #42b983;
  text-align: center;
}

.section {
  margin: 30px 0;
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 8px;
  background: #f9f9f9;
}

h2 {
  margin-top: 0;
  color: #333;
}

.counter-display {
  text-align: center;
  margin: 20px 0;
}

.count {
  font-size: 48px;
  font-weight: bold;
  color: #42b983;
}

.button-group {
  display: flex;
  gap: 10px;
  justify-content: center;
  margin: 20px 0;
}

button {
  padding: 10px 20px;
  font-size: 16px;
  cursor: pointer;
  background: #42b983;
  color: white;
  border: none;
  border-radius: 4px;
  transition: background 0.3s;
}

button:hover {
  background: #35a372;
}

.remove-btn {
  background: #e74c3c;
  padding: 5px 10px;
  font-size: 14px;
  margin-left: 10px;
}

.remove-btn:hover {
  background: #c0392b;
}

.status {
  text-align: center;
  margin-top: 20px;
  padding: 10px;
  border-radius: 4px;
  font-weight: bold;
}

.status.positive {
  background: #d4edda;
  color: #155724;
}

.status.negative {
  background: #f8d7da;
  color: #721c24;
}

.status.ready {
  background: #d1ecf1;
  color: #0c5460;
}

.timestamp {
  text-align: center;
  font-size: 18px;
  color: #666;
}

.input-group {
  display: flex;
  gap: 10px;
  margin-bottom: 15px;
}

input {
  flex: 1;
  padding: 10px;
  font-size: 16px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.items-list {
  list-style: none;
  padding: 0;
}

.items-list li {
  padding: 10px;
  margin: 5px 0;
  background: white;
  border-radius: 4px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.empty-message {
  text-align: center;
  color: #999;
  font-style: italic;
}
</style>
