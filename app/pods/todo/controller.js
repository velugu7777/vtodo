import Ember from 'ember';

export default Ember.Controller.extend({
  actions: {
     createRecord: function(name,isCompleted) {
      this.store.createRecord('todo', {
         name: 'Rails is Omakase',
         isCompleted: false
       });

     }
   }
});
