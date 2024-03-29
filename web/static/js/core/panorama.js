(function (a) {
	a.fn.parallax = function (b) {
		var b = a.extend({useHTML: true, elements: []}, b || {});
		a((b.useHTML) ? "html" : this).mousemove(function (k) {
			var g = a(this);
			var d = {x: Math.floor(parseInt(g.width()) / 2), y: Math.floor(parseInt(g.height()) / 2)};
			var l = {x: (k.pageX - g.offset().left), y: (k.pageY - g.offset().top)};
			var h = {x: (l.x - d.x), y: (l.y - d.y)};
			for (var j = b.elements.length - 1; j >= 0; j--) {
				var c = {}, m, f;
				for (var n in b.elements[j].properties.x) {
					f = b.elements[j].properties.x[n];
					m = f.initial + (h.x * f.multiplier);
					if ("min" in f && m < f.min) {
						m = f.min
					} else {
						if ("max" in f && m > f.max) {
							m = f.max
						}
					}
					if ("invert" in f && f.invert) {
						m = -(m)
					}
					if (!("unit" in f)) {
						f.unit = "px"
					}
					c[n] = m + f.unit
				}
				for (var n in b.elements[j].properties.y) {
					f = b.elements[j].properties.y[n];
					m = f.initial + (h.y * f.multiplier);
					if ("min" in f && m < f.min) {
						m = f.min
					} else {
						if ("max" in f && m > f.max) {
							m = f.max
						}
					}
					if ("invert" in f && f.invert) {
						m = -(m)
					}
					if (!("unit" in f)) {
						f.unit = "px"
					}
					c[n] = m + f.unit
				}
				if ("background-position-x" in c || "background-position-y" in c) {
					c["background-position"] = "" + (("background-position-x" in c) ? c["background-position-x"] : "0px") + " " + (("background-position-y" in c) ? c["background-position-y"] : "0px");
					delete c["background-position-x"];
					delete c["background-position-y"]
				}
				a(b.elements[j].selector).css(c)
			}
		})
	}
})(jQuery);