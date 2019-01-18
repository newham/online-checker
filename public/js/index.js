var app = new Vue({
    delimiters: ['${', '}'],
    el: '#app',
    data: {
        statusList: [{}],
        userAgent: navigator.userAgent.toLowerCase(),
        list_id: -2,
        connect_id: -1,
        refresh_list_time: 1000,
        reconnect_time: 3000,
        btn_connect_show: false,
        userAgentInfo: {
            //操作系统
            os: "",
            browser: ""
        }

    },
    methods: {
        getStatusList: function () {
            this.$http.get('/status/list').then(function (res) {
                // console.log('getStatusList success');
                this.show_disconnect(false);
                this.statusList = res.data;
            }, function (res) {
                console.log('getStatusList error');
                if (res.status == 403) {
                    console.log('unsignined');
                    window.location.href = "/";
                } else {
                    this.statusList = null;
                    this.show_disconnect(true);
                }
            });
        },
        show_disconnect: function (t) {
            if (t) {
                $("#img_disconnect").show();
                this.show_btn_connect(false);
            } else {
                $("#img_disconnect").hide();
            }
        },
        do_getStatusList: function () {
            this.getStatusList();
            this.list_id = setInterval(() => {
                this.getStatusList();
            }, this.refresh_list_time);
        },
        getJson: function () {
            return { userAgent: this.userAgent, os: this.userAgentInfo.os.toLowerCase() };
        },
        connect: function () {
            //发送 post 请求
            this.$http.post('/status/connect', this.getJson()).then(function (res) {
                this.getStatusList();
                // console.log("connect success");
            }, function (res) {
                console.log("connect error");
                // window.clearInterval(this.connect_id);
            });

        },
        disconnect: function () {
            //发送 post 请求
            this.$http.post('/status/disconnect', this.getJson()).then(function (res) {
                this.getStatusList();
                // alert("disconnect success");
                // console.log("disconnect success");
            }, function (res) {
                console.log("disconnect error");
            });

        },
        do_connect: function () {
            this.connect();
            this.connect_id = setInterval(() => {
                this.connect();
            }, this.reconnect_time);
            this.show_btn_connect(false);
        },
        do_disconnect: function () {
            this.disconnect();
            window.clearInterval(this.connect_id);
            this.show_btn_connect(true);
        },
        get_local_time: function (nS) {
            return new Date(parseInt(nS) * 1000).toLocaleString().replace(/:\d{1,2}$/, ' ');
        },
        show_btn_connect: function (show) {
            this.btn_connect_show = show;
        },
        signout: function () {
            this.do_disconnect();
            //发送 post 请求
            this.$http.post('/signout', {}).then(function (res) {

            }, function (res) {
                console.log("signout error");
            });
        },
        parse_user_agent: function () {
            let regs = {};
            let terminal = {
                'windows nt 10': 'Windows 10',
                'windows nt 6.3': 'Windows 8.1',
                'windows nt 6.2': 'Windows 8',
                'windows nt 6.1': 'Windows 7',
                'windows nt 6.0': 'Windows Vista',
                'windows nt 5.2': 'Windows Server 2003XP x64',
                'windows nt 5.1': 'Windows XP',
                'windows xp': 'Windows XP',
                'windows nt 5.0': 'Windows 2000',
                'windows me': 'Windows ME',
                'win98': 'Windows 98',
                'win95': 'Windows 95',
                'win16': 'Windows 3.11',
                'android': 'Android',
                'ubuntu': 'Ubuntu',
                'linux': 'Linux',
                'iphone': 'iPhone',
                'ipod': 'iPod',
                'ipad': 'iPad',
                'macintosh|mac os x': 'Mac OS X',
                'mac_powerpc': 'Mac OS 9',
                'blackberry': 'BlackBerry',
                'webos': 'Mobile',
                'freebsd': 'FreeBSD',
                'sunos': 'Solaris'
            };

            var userAgent = this.userAgent;

            var os = "Others"
            for (let key in terminal) {
                if (new RegExp(key).test(userAgent)) {
                    os = terminal[key];
                    break;
                }
            }


            var browser_name = 'Others';
            if (userAgent.indexOf('firefox') > -1) {
                browser_name = 'Firefox';
            } else if (userAgent.indexOf('chrome') > -1) {
                browser_name = 'Chrome';
            } else if (userAgent.indexOf('trident') > -1 && userAgent.indexOf('rv:11') > -1) {
                browser_name = 'IE11';
            } else if (userAgent.indexOf('msie') > -1 && userAgent.indexOf('trident') > -1) {
                browser_name = 'IE(8-10)';
            } else if (userAgent.indexOf('msie') > -1) {
                browser_name = 'IE(6-7)';
            } else if (userAgent.indexOf('opera') > -1) {
                browser_name = 'Opera';
            } else {
                console.log('Failed to identify the browser,info:' + userAgent);
            }

            if (this.userAgentInfo.os != os) {
                this.userAgentInfo.os = os;
            }
            if (this.userAgentInfo.browser != browser_name) {
                this.userAgentInfo.browser = browser_name;
            }

        }
    },
    mounted() {
        this.parse_user_agent();
        this.do_getStatusList();
        this.do_connect();

    }
});

window.onbeforeunload = function () {
    app.do_disconnect();
}

window.onunload = function () {
    app.do_disconnect();
}

$(function () { $("[data-toggle='tooltip']").tooltip(); });
