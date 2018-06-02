(function () {
  'use strict';

  Vue.component('room', {
    props: {
      data: {
        id: "",
        name: "",
        image_url: "",
        bitrh: "",
      }
    },
    data() {
      return {
        editing: false,
        editedBody: null,
      }
    },

    template: `
    <div class="col">
      <div class="card" style="width: 18rem;">
        <img class="card-img-top" v-bind:src="data.image_url" alt="Card image cap">
          <div class="card-body">
            <h5 class="card-title">{{ data.name }}</h5>
            <p class="card-text">{{ data.birth }}</p>
           <a v-bind:href="'./rooms/' + data.id" class="btn btn-primary">部屋に入る</a>
        </div>
      </div>
    </div>
  `,
  mounted(){
    
  },
    methods: {
      remove() {
        this.removeMessage(this.id)
      }
    }
  });

  const app = new Vue({
    el: '#app',
    data: {
      rooms: [
        {
          id: "1",
          name: "namae",
          image_url: "https://nekogazou.com/wp-content/uploads/2015/03/282e6ed4976b181c78381a609c0f4e32-e1427784795864.jpg",
          bitrh: "",
        },
        {
          id: "2",
          name: "namae",
          image_url: "https://nekogazou.com/wp-content/uploads/2015/03/282e6ed4976b181c78381a609c0f4e32-e1427784795864.jpg",
          bitrh: "",
        },
        {
          id: "3",
          name: "namae",
          image_url: "https://nekogazou.com/wp-content/uploads/2015/03/282e6ed4976b181c78381a609c0f4e32-e1427784795864.jpg",
          bitrh: "",
        },
        {
          id: "4",
          name: "namae",
          image_url: "https://nekogazou.com/wp-content/uploads/2015/03/282e6ed4976b181c78381a609c0f4e32-e1427784795864.jpg",
          bitrh: "",
        },
        {
          id: "5",
          name: "namae",
          image_url: "https://nekogazou.com/wp-content/uploads/2015/03/282e6ed4976b181c78381a609c0f4e32-e1427784795864.jpg",
          bitrh: "",
        },{
          id: "6",
          name: "namae",
          image_url: "https://nekogazou.com/wp-content/uploads/2015/03/282e6ed4976b181c78381a609c0f4e32-e1427784795864.jpg",
          bitrh: "",
        }
      ]
    },
    created() {
      this.getRooms();
    },
    methods: {
       getRooms() {
         fetch('/api/rooms').then(response => response.json()).then(data => {
           this.rooms = data.result;
         });
       }
    }
  });
})();
