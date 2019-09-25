// JavaScript Document
$(document).ready(function(e) {
	  $(".menter_btn_a_a_lf").click(function(){
		if($(".check_boxId").is(":checked")){
			
	     	window.location.href="/install/2";
		}
		else
		{
			alert("请同意安装协议");
		}
	});
});