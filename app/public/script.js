

  // Enable pusher logging - don't include this in production
  Pusher.logToConsole = true;
  $("#list").hide()

  var pusher = new Pusher(PusherKey, {
    cluster: 'eu',
    encrypted: true
  });

  var channel = pusher.subscribe('Domains');

  var w = $("#word")
  var l = $("#list")

  $("#send").click(function() {
    if (w.val()) {
      $.get("/domains?word=" + w.val(), function(response) {
        var key = response
        l.html("")
        w.val("")
        $("#list").hide()
        channel.bind('Result-'+key, function(data) {
          $("#list").show()
          var a = data.available ? "" : "un";
          var html = "<li class='collection-item item-" + a +"available'>"
          html += "<a href='http://" + data.name +"' target='_blank'>" + data.name + "</a>"
          html += " is " + a + "available. </li>"

          l.html(l.html() + html)
        });
        channel.bind('End-'+key, function(data) {
          console.log("Unbind from " + key)
          channel.unbind('Result-'+key)
          channel.unbind('End-'+key)
        });
      })
    }
  })


  function genKey() {
    return "" + Date.now() + "_" + Math.floor(Math.random() * 1000)
  }
