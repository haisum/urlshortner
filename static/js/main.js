
var app = {
	init : function(){
	  $("#shorten-form").attr("action", "javascript:void(0)").on("submit", app.shorten);
	  $("#logout-form").attr("action", "javascript:void(0)").on("submit", app.logout);
	  $("#login-form").attr("action", "javascript:void(0)").on("submit", app.login);
	  $("#register-form").attr("action", "javascript:void(0)").on("submit", app.register);
	  app.templates.shorten = Handlebars.compile($("#shorten-template").html());
	  if($("#errors-template").length > 0){
	  	app.templates.errors = Handlebars.compile($("#errors-template").html());
	  }
	  if($("#data-template").length > 0){
	  	app.templates.data = Handlebars.compile($("#data-template").html());
	  	app.loadData(0, 5);
	  }
	},
	loadData : function(offset, limit){
		$("#data-load").addClass("active");
		$.ajax({
			url : "/me",
			dataType : "json",
			data : {
				offset : offset,
				limit : limit,
			}
		}).done(function(data){
			if(data.Success){
				$("#data-table").html(app.templates.data(data));
				currentPage = (offset/limit +1);
				app.paginate(currentPage, currentPage, data.Total, 5, "#pagination", "#pagination-meta")
			}
			$("#data-load").removeClass("active");
		});
	},
	shorten : function(e){
		$("#shortened-link").html("");
		$("#shorten-load").addClass("active");
		e.preventDefault();
		$.ajax({
			url : "/shorten",
			data : $(e.target).serialize(),
			dataType : "json",
			type : "post"
		}).done(function(data){
			try{
				grecaptcha.reset()
			} catch(err) {

			}
			$("#shorten-load").removeClass("active");
			$("#shortened-link").html(app.templates.shorten(data));
			app.loadData(0, 5);
			$("#copy-url").select();
		});
		return false;
	},
	login : function(e){
		$("#login-errors").html("");
		$("#login-load").addClass("active");
		e.preventDefault();
		$.ajax({
			url : "/login",
			data : $(e.target).serialize(),
			dataType : "json",
			type : "post"
		}).done(function(data){
			if(data.Success){
				window.location.reload();
			}
			$("#login-load").removeClass("active");
			$("#login-errors").html(app.templates.errors(data));
		});
		return false;
	},
	register : function(e){
		$("#register-errors").html("");
		$("#register-load").addClass("active");
		e.preventDefault();
		$.ajax({
			url : "/register",
			data : $(e.target).serialize(),
			dataType : "json",
			type : "post"
		}).done(function(data){
			if(data.Success){
				window.location.reload();
			}
			$("#register-load").removeClass("active");
			$("#register-errors").html(app.templates.errors(data));
		});
		return false;
	},
	logout : function(e){
		$("#logout-load").addClass("active");
		e.preventDefault();
		$.ajax({
			url : "/logout",
			type : "post",
		}).done(function(){
			window.location.reload();
		});
		return false;
	},
	templates : {},
	paginate : function ( firstli, currentli, totalrecordcount, pagesize, pagination_links_container, pagination_meta_container ) {
		//	var pagenumber;
		var max_pages = 5;	
		// determin the last page: lastpage
		var lastpage = Math.ceil(totalrecordcount / pagesize);
		// determin the first page: firstli
		if (firstli == currentli )  firstli = currentli - parseInt(max_pages/2);
		if ( firstli < 1 ) firstli = 1;
		if ( ( lastpage - firstli ) < max_pages && lastpage > (max_pages-1) ) firstli = lastpage - (max_pages-1);
		if ( lastpage < (max_pages+1) ) firstli = 1;

		// determin the currentli
		//if ( currentli < max_pages ) currentli = '0' + currentli; else currentli = currentli;


		/* first step, reset the navigation tool */
		var html_li = ""
		$(pagination_links_container).html("");
		$(pagination_meta_container).html("");
		if ( firstli > 1 && lastpage >= max_pages ) html_li = html_li + '<li class="waves-effect"><a href="javascript:;">1</a></li>'; 
		if ( lastpage > max_pages && firstli > 1 ) html_li = html_li + '<li class="waves-effect"><a href="javascript:;" class="first"><i class="mdi-navigation-chevron-left"></i></a></li>'; 

		for ( var i = 0 ; i < max_pages;  i++ ) 
			{
			var sum = i+firstli;
			if ( sum > lastpage ) break;
			if ( i==0 ) 
				{ 
				if ( currentli == sum )
					html_li = html_li + '<li class="firstli currentli waves-effect active"><a href="javascript:;">' + sum + "</a></li>";
				else
					html_li = html_li + '<li class="firstli waves-effect"><a href="javascript:;">' + sum + "</a></li>";
				} 
			else if (i==(max_pages-1)) 
				{ 
				if ( currentli == sum )
					html_li = html_li + '<li class="lastli currentli active waves-effect"><a href="javascript:;">' + sum + "</a></li>"; 
				else
					html_li = html_li + '<li class="lastli waves-effect"><a href="javascript:;">' + sum + "</a></li>"; 
				} 
			else 
				{ 
				if ( currentli == sum )
					html_li = html_li + '<li class="currentli active waves-effect"><a href="javascript:;">' + sum + "</a></li>"; 
				else
					html_li = html_li + '<li class="waves-effect"><a href="javascript:;">' + sum + "</a></li>"; 
				}
			};

		if ( lastpage > max_pages && (lastpage - firstli) > max_pages ) html_li = html_li + '<li class="waves-effect"><a class="last" href="javascript:;"><i class="mdi-navigation-chevron-right"></a></li>';
		if ( lastpage > max_pages && (lastpage - firstli) >= max_pages ) html_li = html_li + '<li  class="waves-effect"><a href="javascript:;">' + lastpage + '</a></li>';

		/* last step, write the navigation tool */ 
		$(pagination_links_container).html(html_li);
		var recordstart = ( pagesize * (currentli - 1 ) ) + 1;
		if  ( totalrecordcount > (recordstart + pagesize) ) var recordend = parseInt(recordstart) + parseInt(pagesize) - 1;
		else var recordend = totalrecordcount;
		if(recordend > totalrecordcount){
			recordend = totalrecordcount;
		}
		if ( totalrecordcount > 2 ) $(pagination_meta_container).html( "Records " + recordstart + " - " + recordend + " of " + totalrecordcount );//Records 1 - 25 of max_pages

		/* add the click event */
		$(pagination_links_container + ">li").click( function (event) {
			pagenumber = $(this).text(); 
			firstli = Number( $(pagination_links_container + ">.firstli:first>a").text() ) ;
			lastli = Number( $(pagination_links_container + ">lastli:first>a").text() );

			if ( $(this).find("a").hasClass("first") ) 
				{
				firstli = parseInt(firstli) + parseInt(max_pages);
				app.paginate( firstli, currentli, totalrecordcount, pagesize, pagination_links_container,pagination_meta_container);
				}
			if ( $(this).find("a").hasClass("last") ) 
				{
				firstli = firstli - max_pages;
				if ( firstli < 1 ) firstli = 1;
				app.paginate( firstli, currentli, totalrecordcount, pagesize, pagination_links_container,pagination_meta_container);
				}
			if ( Number(pagenumber) > 0 ) {
				pagenumber = Number(pagenumber);
				app.loadData((pagenumber-1) * 5, 5);
			};
		});
	}
}


$(app.init);