import Ember from 'ember';

export default Ember.Controller.extend({
  actions: {
     creteTodo: function(post) {
       console.log(post.get('name'));
     }
   }
});
