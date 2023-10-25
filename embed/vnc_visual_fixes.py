import os
import re

title = os.environ.get("HTTP_TITLE", "Sign in to your account")

replace = """<title>[TITLE]</title>
<link rel=icon type=image/png href=app/images/icons/favicon.png>
<style>
#noVNC_transition, #noVNC_control_bar, #noVNC_status, #noVNC_control_bar_handle {
    display: none !important;
}
</style>
<script>
try {
    window.originalTitle = document.title; // save for future
    Object.defineProperty(document, 'title', {
        get: function() {return originalTitle},
        set: function() {}
    });
} catch (e) {}
</script>"""

replace = replace.replace("[TITLE]", title)

with open("/usr/share/kasmvnc/www/vnc.html") as f:
    contents = f.read()
    contents = re.sub(r"<link rel=icon [^<]+>", "", contents)
    contents = contents.replace("<title>KasmVNC</title>", replace)

with open("/usr/share/kasmvnc/www/vnc.html", "w") as f:
    f.write(contents)
