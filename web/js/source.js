(function (global) {
  var COOKIE = "ldy_source";
  var PREFIX = "channel://";
  var SOURCE_RE = /^[A-Za-z0-9][A-Za-z0-9_-]{1,31}$/;
  var RESERVED = {
    official: 1,
    "404": 1,
    invite: 1,
    source: 1,
    channel: 1,
    share: 1,
    agent: 1,
  };

  function parseSource(raw) {
    var s = String(raw || "").trim();
    if (!s) return "";
    if (s.indexOf("share://") >= 0 || s.indexOf("agent://") >= 0) return "";
    var i = s.indexOf(PREFIX);
    if (i >= 0) s = s.slice(i + PREFIX.length).trim();
    s = s.split(/[?&#/\s]/)[0] || "";
    if (!SOURCE_RE.test(s) || RESERVED[s.toLowerCase()]) return "";
    return s;
  }

  function sourceFromLocation(loc) {
    loc = loc || window.location;
    var params = new URLSearchParams(loc.search || "");
    var fromQuery = parseSource(params.get("source"));
    if (fromQuery) return fromQuery;
    var path = String(loc.pathname || "");
    var m = path.match(/\/s\/([^/]+)/i);
    if (m) return parseSource(decodeURIComponent(m[1]));
    return "";
  }

  function readCookie() {
    var parts = ("; " + document.cookie).split("; " + COOKIE + "=");
    if (parts.length < 2) return "";
    return parseSource(decodeURIComponent(parts.pop().split(";").shift()));
  }

  function writeCookie(code) {
    if (!code) return;
    var maxAge = 5 * 24 * 3600;
    document.cookie = COOKIE + "=" + encodeURIComponent(code) + ";path=/;max-age=" + maxAge + ";SameSite=Lax";
  }

  function rememberSource(loc) {
    var code = sourceFromLocation(loc) || readCookie();
    if (code) writeCookie(code);
    return code;
  }

  function clipboardValue(code) {
    return code ? PREFIX + code : "";
  }

  global.LdySource = {
    parseSource: parseSource,
    rememberSource: rememberSource,
    clipboardValue: clipboardValue,
  };
})(window);
