package main

/* HTML template we'll render over in main.go's HTTP handler.

It expects the following input:

WisdomOutput    string    Quote to display
DebugOutput     string    Debug output to display
VersionOutput   string    Version and build info to display
ShowDebug       bool      Whether to expand debug output by default

*/

//nolint:lll
// ... since I don't mind long lines in my HTML templates.
const html = `
<!doctype html>
<html>
<head>
<meta name="viewport" content="width=device-width, initial-scale=1" charset="utf-8">
<style type="text/css">

html, body, p {
  font-family: sans-serif;
}

.collapsible {
  cursor: pointer;
  border: yes;
  outline: none;
  font-size: 14px;
}

.collapsible:hover {
  background-color: light-grey;
}

.active {
  background-color: light-grey;
}

.collapsible:after {
  content: '+';
  float: right;
  margin-left: 5px;
}

.active:after {
  content: "-";
}

.debug {
  max-height: 0;
  padding: 0 12px;
  overflow: hidden;
  transition: max-height 0.1s ease-out;
  background-color: #f1f1f1;
  background-color: #000000;
  color: #00ff00;
  overflow-wrap: break-word;
}

.quote {
  text-align: left;
  margin: 0 auto;
  max-width: 80ch;
  white-space: pre-wrap;
  //word-break: normal;
  //overflow-wrap: break-word;
  //border: 1px solid red;
}

.reload-button {
  font-size: 50px;
  text-align: center;
  text-decoration: none;
  //border: 1px solid blue;
}
</style>
<title>🥠 An adage a day keeps the boredom away!</title>
</head>
<body>

<div class="quote"><pre class="quote">
{{.WisdomOutput}}
</pre>
<div class="reload-button"><a class="reload-button" title="Another quote, please." href="javascript:history.go(0)">🎱</a></div></div>
<p>Available backends: {{/* FIXME: Auto-generate backend list and discover options. Go's template package strips all HTML comments, so I didn't have to write this comment in a template action to hide it from the rendered response. */}}
<ul>
  <li><a href="?backend=flatfile">flatfile</a></li>
  <li><a href="?backend=fortune">fortune</a></li>
</ul>
</p>

<hr>
  <h1>Aphorismophilia, <i>n</i>.:</h1>
  <ol>
    <li>Love of aphorisms.</li>
    <li>Flimsy pretense to gain hands-on experience with modern technology trends.</li>
  </ol>
  <p>Learn more on <a href="https://github.com/mikeroach/aphorismophilia">Github</a>.</p>
  
<button class="collapsible">🧘‍♂️ Guru Meditiation:</button>
<div class="debug">
<pre>
{{.DebugOutput}}
{{.VersionOutput}}
</pre>
</div>
<script>
// Adapted from https://www.w3schools.com/howto/howto_js_collapsible.asp
var coll = document.getElementsByClassName("collapsible");
var i;

for (i = 0; i < coll.length; i++) {

  {{ if .ShowDebug }}
  // Display expanded debug output on page load
  coll[i].classList.toggle("active");
  var content = coll[i].nextElementSibling;
  content.style.maxHeight = content.scrollHeight + "px";
  {{- end}}

  coll[i].addEventListener("click", function() {
    this.classList.toggle("active");
    var content = this.nextElementSibling;
    if (content.style.maxHeight){
      content.style.maxHeight = null;
    } else {
      content.style.maxHeight = content.scrollHeight + "px";
    } 
  });
}
</script>
</body>
</html>
`
