
var app = {
	init : function(){
	  $("#copy-url").on("mouseenter", function(){ 
	  	$(this).select()
	  });
	}
}


$(app.init);