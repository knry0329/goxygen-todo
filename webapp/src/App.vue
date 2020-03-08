<template>
  <div>
    <h3>My Todo</h3>
    <input v-model="newTodo" placeholder="input here..."> 
    <button @click="addTodo()">ADD</button>

    <h5>Todo List</h5>
    <ul>
      <li v-for="todo in todos" :key="todo.id">
        <div :id="todo.id">
        {{ todo.text }}
        <button @click="deleteTodo(todo.id)">DEL</button>
      </div>
    </li>
    </ul>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  data() {
    return {
      todos: [],
      newTodo: ''
    }
  },
  mounted() {
    this.fetchTodo()
  },
  methods: {
    fetchTodo() {
      axios
      .get(`${process.env.VUE_APP_API_URL}/api/todos`)
      .then(response => (this.todos = response.data))
    },
    addTodo() {
      if (this.newTodo==='') return;
      axios
        .post(`${process.env.VUE_APP_API_URL}/api/todos`, {
          text: this.newTodo
        })
        .then(() => this.fetchTodo())
      this.newTodo = ''
    },
    deleteTodo(i) {
      axios
        .delete(`${process.env.VUE_APP_API_URL}/api/todos?id=${i}`)
        .then(() => this.fetchTodo())
    }
  }
}
</script>