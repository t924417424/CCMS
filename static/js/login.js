var xsrf, xsrflist;
xsrf = $.cookie("_xsrf");
xsrflist = xsrf.split("|");
_xsrf = Base64.decode(xsrflist[0]);
function userlogin() {
    var username = $("#username").val();
    $.post("/login",{username:username,_xsrf:_xsrf},function(result){
        $('#msg').html(result.msg);
        if(result.code == 1){
            window.location.href = "/";
        }
    });
}
function adminlogin() {
    var username = $("#username").val();
    var password = $("#password").val();
    var captcha = $("#captcha").val();
    var captchas = $("input[name='captchas']")[0].value
    $.post("/ccms/c_login",{username:username,password:password,captcha:captcha,captchas:captchas,_xsrf:_xsrf},function(result){
        $('#msg').html(result.msg);
        if(result.code == 1){
            window.location.href = "/ccms/index";
        }
    });
}
function fileget() {
    $('input[type="button"]').prop('disabled', true);
    var uploadFile = $('#filename')[0].files[0];
    var formdata = new FormData();
    formdata.append('fileInfo', uploadFile);
    formdata.append('_xsrf', _xsrf);
    $.ajax({
        url: '/ccms/upload',
        type: 'post',
        data : formdata,
        processData: false,
        contentType: false,
        xhr: function() {
            var xhr = new XMLHttpRequest();
            //使用XMLHttpRequest.upload监听上传过程，注册progress事件，打印回调函数中的event事件
            xhr.upload.addEventListener('progress', function (e) {
                //console.log(e);
                //loaded代表上传了多少
                //total代表总数为多少
                var progressRate = (e.loaded / e.total) * 100 + '%';
                console.log(progressRate);
                //通过设置进度条的宽度达到效果
                $('#msg').html(progressRate);
            })

            return xhr;
        },
        success: function (data,status) {
            $('#msg').html("解析中");
            if (data.code == 200){
                window.location.href = "/ccms/getexcel";
            }else if (data.code == 100){
                $('#msg').html("文件格式错误！");
            }else{
                $('#msg').html("上传失败！");
            }
        }
    })
}
function sysset() {
    var sitename = $("#sitename").val();
    $.post("/ccms/sysset",{sitename:sitename,_xsrf:_xsrf},function(result){
        $('#msg').html(result.msg);
    });
}
function mysqlconn() {
    var xsrf, xsrflist;
    xsrf = $.cookie("_xsrf");
    xsrflist = xsrf.split("|");
    _xsrf = Base64.decode(xsrflist[0]);
    var dbhost = $("#dbhost").val();
    var dbuser = $("#dbuser").val();
    var dbpas = $("#dbpas").val();
    var dbname = $("#dbname").val();
    var website = $("#website").val();
    var username = $("#adminuser").val();
    var password = $("#adminpas").val();
    var salt = $("#salt").val()
    if(dbhost && dbuser && dbpas && dbname && website && username && password && salt){
        $.post("/install/conn",{dbhost:dbhost,dbuser:dbuser,dbpass:dbpas,dbname:dbname,website:website,username:username,password:password,salt:salt,_xsrf:_xsrf},function(result){
            alert(result.msg);
            if(result.code == 1){
                window.location.href = "/install/4";
            }
        });
    }else{
        alert("数据不能为空！");
    }
}
function setexcel(row) {
    $('#msg').html("正在分析数据源内容。。。");
    $('input[type="button"]').prop('disabled', true);
    $.post("/ccms/importexcel",{row:row,_xsrf:_xsrf},function(result){
        $('#msg').html(result.msg);
        if(result.code == 1){
            window.location.href = "/ccms/datasys";
        }
    });
}
function addtitle(title) {
    var id = "c" + title;
    $("#" + id).attr("disabled",true);
    var values = $('#etitle').val();
    if(values != ""){
        title = values + "," + title
    }
    $('#etitle').val(title);
}
function setdata() {
    var etitle = $('#etitle').val();
    var onlydata = $('#onlytitle').find("option:selected").text();
    var onlypass = $('#passwordtitle').find("option:selected").text();
    if(etitle==""||onlydata==""||onlypass==""||etitle== undefined||onlydata== undefined||onlypass== undefined){
        alert("设置内容不能为空！")
    }else{
        $.post("/ccms/datasys",{etitle:etitle,onlydata:onlydata,onlypass:onlypass,_xsrf:_xsrf},function(result){
            alert(result.msg)
            //window.location.reload();
            /*if(result.code == 1){
                window.location.href = "/";
            }*/
        });
    }
}
function setpass() {
    var oldpass = $('#oldpass').val();
    var newpass = $('#newpass').val();
    var repass = $('#repass').val();
    $.post("/ccms/repass",{oldpass:oldpass,newpass:newpass,repass:repass,_xsrf:_xsrf},function(result){
        alert(result.msg)
        /*if(result.code == 1){
            window.location.href = "/";
        }*/
    });
}
function getselect() {
    var select  = $('#select').val();
    if (select == "" || select == undefined){
        alert("检索内容不能为空！");
    }else{
        var newWeb=window.open('_blank');
        newWeb.location='api/' + select;
    }
}
function getbind() {
    var select  = $('#select').val();
    if (select == "" || select == undefined){
        alert("检索内容不能为空！");
    }else{
        var newWeb=window.open('_blank');
        newWeb.location='bind/' + select;
    }
}
function cbind() {
    var mymessage=confirm("清理绑定关系，用户再次登陆将提示重新绑定，确认？");
    if(mymessage==true){
        $.post("/ccms/bind",{_xsrf:_xsrf},function(result){
            alert(result.msg)
        });
    }
}
function userlogin() {
    var username = $("#user").val();
    var password = $("#pass").val();
    var captcha = $("#captcha").val();
    var captchas = $("input[name='captchas']")[0].value;
    if (username == "" || password =="" || captcha ==""){
        $('#msg').html("输入内容不能为空！");
    }else{
        $.post("/login",{username:username,password:password,captcha:captcha,captchas:captchas,_xsrf:_xsrf},function(result){
            $('#msg').html(result.msg);
            if(result.code == 200){
                var data = JSON.parse(result.data);
                for (var key in data) {
                    $("#data_table tbody").append('<tr><td>' + key + '</td><td>' + data[key] + '</td></tr>');
                    console.log(key);     //获取key值
                    console.log(data[key]); //获取对应的value值
                }
                $('#loginform').animate({
                    height: 'toggle',
                    opacity: 'toggle'
                }, 'slow');
                setInterval(function(){
                    $('#datas').show("slow");
                },1000);

            }
        });
    }
}
//https://graph.qq.com/oauth2.0/authorize?client_id=101488968&response_type=code&redirect_uri={授权回调地址}&state={原样返回的参数}