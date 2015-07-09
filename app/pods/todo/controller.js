import Ember from 'ember';

export default Ember.Controller.extend({
  actions: {
     createRecord: function() {
      this.store.createRecord('todo', {
         name: 'Rails is Omakase',
         isCompleted: false
       });

     }
   }
});
