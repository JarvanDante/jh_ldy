(function () {
  var cfg = {
    site_name: "JH",
    slogan: "打开 App，马上开始",
    android_url: "",
    ios_url: "",
    h5_url: "",
    logo_text: "JH",
    skin_css: "/skins/default/app.css?v=20261026",
    theme_color: "#0b0b0d",
  };
  var source = "";

  function toast(text) {
    var el = document.getElementById("toast");
    el.textContent = text;
    el.hidden = false;
    clearTimeout(toast._t);
    toast._t = setTimeout(function () {
      el.hidden = true;
    }, 1800);
  }

  function inWeixinOrQQ() {
    var ua = navigator.userAgent || "";
    return /MicroMessenger|QQ\//i.test(ua);
  }

  function writeClipboard(text) {
    if (!text) return Promise.resolve(false);
    if (navigator.clipboard && window.isSecureContext) {
      return navigator.clipboard.writeText(text).then(
        function () {
          return true;
        },
        function () {
          return fallbackCopy(text);
        }
      );
    }
    return Promise.resolve(fallbackCopy(text));
  }

  function fallbackCopy(text) {
    var ta = document.createElement("textarea");
    ta.value = text;
    ta.setAttribute("readonly", "");
    ta.style.position = "fixed";
    ta.style.left = "-9999px";
    document.body.appendChild(ta);
    ta.select();
    var ok = false;
    try {
      ok = document.execCommand("copy");
    } catch (e) {
      ok = false;
    }
    document.body.removeChild(ta);
    return ok;
  }

  function bindSourceThenGo(rawUrl, kind) {
    if (inWeixinOrQQ()) {
      document.getElementById("browserMask").hidden = false;
      return;
    }
    if (!rawUrl) {
      toast(kind + "下载地址未配置");
      return;
    }
    writeClipboard(window.LdySource.clipboardValue(source)).then(function () {
      window.location.href = rawUrl;
    });
  }

  function render() {
    document.title = cfg.site_name + " · 下载";
    document.getElementById("siteName").textContent = cfg.site_name;
    document.getElementById("slogan").textContent = cfg.slogan;
    document.getElementById("logo").textContent = cfg.logo_text || "JH";
    if (cfg.skin_css) document.getElementById("skinCss").href = cfg.skin_css;
    if (cfg.theme_color) {
      document.getElementById("themeColor").setAttribute("content", cfg.theme_color);
    }
  }

  function bind() {
    document.getElementById("btnAndroid").addEventListener("click", function () {
      bindSourceThenGo(cfg.android_url, "安卓");
    });
    document.getElementById("btnIos").addEventListener("click", function () {
      bindSourceThenGo(cfg.ios_url, "iOS");
    });
    document.getElementById("browserMask").addEventListener("click", function () {
      this.hidden = true;
    });
  }

  source = window.LdySource.rememberSource();
  bind();
  var boot = "/api/boot";
  if (source) boot += "?source=" + encodeURIComponent(source);
  fetch(boot, { cache: "no-store" })
    .then(function (res) {
      return res.ok ? res.json() : {};
    })
    .catch(function () {
      return {};
    })
    .then(function (data) {
      cfg.site_name = data.site_name || cfg.site_name;
      cfg.slogan = data.slogan || cfg.slogan;
      cfg.android_url = String(data.android_url || "").trim();
      cfg.ios_url = String(data.ios_url || "").trim();
      cfg.h5_url = String(data.h5_url || "").trim();
      cfg.logo_text = data.logo_text || cfg.logo_text;
      cfg.skin_css = data.skin_css || cfg.skin_css;
      cfg.theme_color = data.theme_color || cfg.theme_color;
      if (data.source) source = data.source;
      render();
      writeClipboard(window.LdySource.clipboardValue(source));
    });
})();
