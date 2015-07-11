import Ember from 'ember';

export default Ember.Controller.extend({

  actions: {
    newTodo(name){
      this.store.createRecord('todo',{
        name: name,
        isCompleted: false,
      }).save();
    },

    editTodo(todo){
      todo.set('isVisible', true);
    },

   }

});
