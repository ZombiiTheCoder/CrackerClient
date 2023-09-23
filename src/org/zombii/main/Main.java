package org.zombii.main;

import java.net.Inet4Address;

import ca.weblite.webview.WebView;
import ca.weblite.webview.WebViewCLI;
import ca.weblite.webview.WebViewCLIClient;
import ca.weblite.webview.WebviewSocketServer;

public class Main {

    public static void main(String[] args) throws Exception {
        if (WebViewCLI.restartJVM(args)) {
            return;
        }
        System.out.println(Inet4Address.getLocalHost().getHostAddress());
        WebView w2 = new WebView();
        WebViewCLIClient w = new WebViewCLIClient(
                new String[] { "http://" + Inet4Address.getLocalHost().getHostAddress() + ":2044" });
        w2.url("http://192.168.119.74:2044/");
        w2.show();
        w.close();

    }

}