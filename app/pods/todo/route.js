import Ember from 'ember';

export default Ember.Route.extend({
  model(){
    return this.store.find('todo');
  },
  actions: {
    delTodo(id) {
      this.store.find('todo', id).then(function (todo) {
        todo.destroyRecord(); // => DELETE to /todos/id
      });
    }
  }
});
