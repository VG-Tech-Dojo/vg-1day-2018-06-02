(function() {
  'use strict';
  const Message = function() {
    this.body = '';
    this.username = '';
    this.userid = '';
    this.type = '2';
  };
  const User = function() {
    this.id = '';
    this.name = '';
    this.point = '';
  }

  Vue.component('message', {
    // Tutorial 1-2. ユーザー名を表示しよう
    props: ['id', 'body', 'username', 'removeMessage', 'updateMessage'],
    data() {
      return {
        editing: false,
        editedBody: null,
      }
    },
    // Tutorial 1-2. ユーザー名を表示しよう
    template: `
    <div class="message">
      <div v-if="editing">
        <div class="row">
          <textarea v-model="editedBody" class="u-full-width"></textarea>
        </div>
      </div>
      <div class="message-body" v-else>
        <span>{{ body }} - {{ username }}</span>
      </div>
    </div>
  `,
    methods: {
      remove() {
        this.removeMessage(this.id)
      },
      edit() {
        this.editing = true
        this.editedBody = this.body
      },
      cancelEdit() {
        this.editing = false
        this.editedBody = null
      },
      doneEdit() {
        this.updateMessage({id: this.id, body: this.editedBody})
          .then(response => {
            this.cancelEdit()
          })
      }
    }
  });
  Vue.component('user', {
    props: ['id', 'point', 'name'],
    data() {
      return {
      }
    },
    template: `
    <h4>{{ name }}</h4>
  `
  });

  const app = new Vue({
    el: '#app',
    data: {
      messages: [],
      newMessage: new Message(),
      user: undefined,
    },
    created() {
      this.getMessages();
      this.getUserInfo();
    },
    methods: {
      getMessages() {
        fetch('/api/messages').then(response => response.json()).then(data => {
          this.messages = data.result;
        });
      },
      getUserInfo(){
        fetch('/api/user').then(response => response.json()).then(data => {
          this.user = {
            id: 1,
            point: 0,
            name: 'TEST_NAME'
          };
          // this.user = data.result;
        });
      },
      sendMessage() {
        const message = this.newMessage;
        fetch('/api/messages', {
          method: 'POST',
          body: JSON.stringify(message)
        })
          .then(response => response.json())
          .then(response => {
            if (response.error) {
              alert(response.error.message);
              return;
            }
            this.messages.push(response.result);
            this.clearMessage();
          })
          .catch(error => {
            console.log(error);
          });
      },
      removeMessage(id) {
        return fetch(`/api/messages/${id}`, {
          method: 'DELETE'
        })
        .then(response => response.json())
        .then(response => {
          if (response.error) {
            alert(response.error.message);
            return;
          }
          this.messages = this.messages.filter(m => {
            return m.id !== id
          })
        })
      },
      updateMessage(updatedMessage) {
        return fetch(`/api/messages/${updatedMessage.id}`, {
          method: 'PUT',
          body: JSON.stringify(updatedMessage),
        })
        .then(response => response.json())
        .then(response => {
            if (response.error) {
              alert(response.error.message);
              return;
            }
            const index = this.messages.findIndex(m => {
              return m.id === updatedMessage.id
            })
            Vue.set(this.messages, index, response.result)
        })
      },
      clearMessage() {
        this.newMessage = new Message();
      }
    }
  });
})();
