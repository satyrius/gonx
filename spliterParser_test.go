package gonx

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func BenchmarkSpliterParser(b *testing.B) {
	input := "183.202.32.72|#|-|#|[08/Aug/2023:08:25:35 +0800]|#|1691454335.995|#|release-qn.233lyly.com|#|https|#|Hit|#|GET /upload/spread/generated/31931dd328e65ad83de123559af000fa1bb30253/233\xE4\xB9\x90\xE5\x9B\xAD.apk HTTP/1.1|#|206|#|788105|#|bytes=24379392-25165823|#|bytes 24379392-25165823/26096807|#|-|#|Dalvik/2.1.0 (Linux; U; Android 11; V2156A Build/RP1A.200720.012)|#|-|#|0.184|#|0.171|#|0.004|#|0.004|#|cache27.l2nu20-1[51,51,206-0,M], cache11.l2nu20-1[52,0], cache11.l2nu20-1[57,0], cache1.cn813[99,99,206-0,M], cache2.cn813[117,0], xh-cn3525[217,216,200-0,M], xh-cn3525[220,0], [9217,edge-hb-wuhan13-cache-37.in.ctcdn.cn], [2,edge-ha-zhengzhou29-cache-38.in.ctcdn.cn]|#|127.0.0.1:8088|#|application/vnd.android.package-archive|#|QNM:ac2f4c81149d8b55d22d803ef925b872;QNM3|#|R5dbKgPXi-kLj4RGdGQ|#|ac2f4c81149d8b55d22d803ef925b872|#|117.161.239.165|#|22443|#|4|#|786432|#|368|#|786432|#|788109|#|206|#|-|#|100062|#|6534|#|104|#|14600|#|2|#|-|#|-|#|qiniudcdn|#|-|#|-|#|0|#|0|#|5|#|184|#|qnvm|#|remote_from=client|#|Qypid=-|#|Qyid=-|#|request_completion=OK|#|c_auth_time=-|#|-|#|5115|#|redirect_to=-|#|-|#|-|#|-|#|-|#|-|#|Mon, 07 Aug 2023 13:52:57 GMT|#|117.161.239.165|#|options=-"

	pattern := "$real_remote_ip|#|$remote_user|#|[$time_local]|#|$msec|#|$real_host|#|$scheme|#|$x_qnm_cache|#|$request_method $real_request_uri $server_protocol|#|$status|#|$bytes_sent|#|$http_range|#|$sent_http_content_range|#|$http_referer|#|$http_user_agent|#|$http_x_forwarded_for|#|$request_time|#|$upstream_response_time|#|$upstream_connect_time|#|$upstream_header_time|#|$upstream_http_via|#|$upstream_addr|#|$upstream_http_content_type|#|$upstream_http_x_m_log|#|$master_reqid-$child_reqid|#|$real_hostname|#|$real_server_addr|#|$server_port|#|$first_byte_time|#|$sent_http_content_length|#|$request_length|#|$body_bytes_sent|#|$upstream_bytes_received|#|$upstream_status|#|$cache_level|#|$tcpinfo_rtt|#|$tcpinfo_rttvar|#|$tcpinfo_snd_cwnd|#|$tcpinfo_rcv_space|#|$connection_requests|#|$category|#|$internal_code|#|$http_x_from_cdn|#|$http_x_from_fsrcproxy|#|$response_from|#|$access_time|#|$access_end_time|#|$first_body_bytes_time|#|$log_time|#|$idc_name|#|remote_from=$remote_from|#|Qypid=$qypid|#|Qyid=$qyid|#|request_completion=$request_completion|#|c_auth_time=$c_auth_time|#|$http_content_type|#|$remote_port|#|redirect_to=$redirect_to|#|$cache_status"

	p := NewSpliteParser(pattern, "|#|")
	for i := 0; i < b.N; i++ {
		p.ParseString(input)
	}
}

func BenchmarkStdParser(b *testing.B) {
	input := "183.202.32.72|#|-|#|[08/Aug/2023:08:25:35 +0800]|#|1691454335.995|#|release-qn.233lyly.com|#|https|#|Hit|#|GET /upload/spread/generated/31931dd328e65ad83de123559af000fa1bb30253/233\xE4\xB9\x90\xE5\x9B\xAD.apk HTTP/1.1|#|206|#|788105|#|bytes=24379392-25165823|#|bytes 24379392-25165823/26096807|#|-|#|Dalvik/2.1.0 (Linux; U; Android 11; V2156A Build/RP1A.200720.012)|#|-|#|0.184|#|0.171|#|0.004|#|0.004|#|cache27.l2nu20-1[51,51,206-0,M], cache11.l2nu20-1[52,0], cache11.l2nu20-1[57,0], cache1.cn813[99,99,206-0,M], cache2.cn813[117,0], xh-cn3525[217,216,200-0,M], xh-cn3525[220,0], [9217,edge-hb-wuhan13-cache-37.in.ctcdn.cn], [2,edge-ha-zhengzhou29-cache-38.in.ctcdn.cn]|#|127.0.0.1:8088|#|application/vnd.android.package-archive|#|QNM:ac2f4c81149d8b55d22d803ef925b872;QNM3|#|R5dbKgPXi-kLj4RGdGQ|#|ac2f4c81149d8b55d22d803ef925b872|#|117.161.239.165|#|22443|#|4|#|786432|#|368|#|786432|#|788109|#|206|#|-|#|100062|#|6534|#|104|#|14600|#|2|#|-|#|-|#|qiniudcdn|#|-|#|-|#|0|#|0|#|5|#|184|#|qnvm|#|remote_from=client|#|Qypid=-|#|Qyid=-|#|request_completion=OK|#|c_auth_time=-|#|-|#|5115|#|redirect_to=-|#|-|#|-|#|-|#|-|#|-|#|Mon, 07 Aug 2023 13:52:57 GMT|#|117.161.239.165|#|options=-"

	pattern := "$real_remote_ip|#|$remote_user|#|[$time_local]|#|$msec|#|$real_host|#|$scheme|#|$x_qnm_cache|#|$request_method $real_request_uri $server_protocol|#|$status|#|$bytes_sent|#|$http_range|#|$sent_http_content_range|#|$http_referer|#|$http_user_agent|#|$http_x_forwarded_for|#|$request_time|#|$upstream_response_time|#|$upstream_connect_time|#|$upstream_header_time|#|$upstream_http_via|#|$upstream_addr|#|$upstream_http_content_type|#|$upstream_http_x_m_log|#|$master_reqid-$child_reqid|#|$real_hostname|#|$real_server_addr|#|$server_port|#|$first_byte_time|#|$sent_http_content_length|#|$request_length|#|$body_bytes_sent|#|$upstream_bytes_received|#|$upstream_status|#|$cache_level|#|$tcpinfo_rtt|#|$tcpinfo_rttvar|#|$tcpinfo_snd_cwnd|#|$tcpinfo_rcv_space|#|$connection_requests|#|$category|#|$internal_code|#|$http_x_from_cdn|#|$http_x_from_fsrcproxy|#|$response_from|#|$access_time|#|$access_end_time|#|$first_body_bytes_time|#|$log_time|#|$idc_name|#|remote_from=$remote_from|#|Qypid=$qypid|#|Qyid=$qyid|#|request_completion=$request_completion|#|c_auth_time=$c_auth_time|#|$http_content_type|#|$remote_port|#|redirect_to=$redirect_to|#|$cache_status"

	p := NewParser(pattern)
	for i := 0; i < b.N; i++ {
		p.ParseString(input)
	}
}

func TestSpliterParser(t *testing.T) {
	convey.Convey("Nginx format parser with spliter", t, func() {
		input := "183.202.32.72|#|-|#|[08/Aug/2023:08:25:35 +0800]|#|1691454335.995|#|release-qn.233lyly.com|#|https|#|Hit|#|GET /upload/spread/generated/31931dd328e65ad83de123559af000fa1bb30253/233\xE4\xB9\x90\xE5\x9B\xAD.apk HTTP/1.1|#|206|#|788105|#|bytes=24379392-25165823|#|bytes 24379392-25165823/26096807|#|-|#|Dalvik/2.1.0 (Linux; U; Android 11; V2156A Build/RP1A.200720.012)|#|-|#|0.184|#|0.171|#|0.004|#|0.004|#|cache27.l2nu20-1[51,51,206-0,M], cache11.l2nu20-1[52,0], cache11.l2nu20-1[57,0], cache1.cn813[99,99,206-0,M], cache2.cn813[117,0], xh-cn3525[217,216,200-0,M], xh-cn3525[220,0], [9217,edge-hb-wuhan13-cache-37.in.ctcdn.cn], [2,edge-ha-zhengzhou29-cache-38.in.ctcdn.cn]|#|127.0.0.1:8088|#|application/vnd.android.package-archive|#|QNM:ac2f4c81149d8b55d22d803ef925b872;QNM3|#|R5dbKgPXi-kLj4RGdGQ|#|ac2f4c81149d8b55d22d803ef925b872|#|117.161.239.165|#|22443|#|4|#|786432|#|368|#|786432|#|788109|#|206|#|-|#|100062|#|6534|#|104|#|14600|#|2|#|-|#|-|#|qiniudcdn|#|-|#|-|#|0|#|0|#|5|#|184|#|qnvm|#|remote_from=client|#|Qypid=-|#|Qyid=-|#|request_completion=OK|#|c_auth_time=-|#|-|#|5115|#|redirect_to=-|#|-|#|-|#|-|#|-|#|-|#|Mon, 07 Aug 2023 13:52:57 GMT|#|117.161.239.165|#|options=-"

		pattern := "$real_remote_ip|#|$remote_user|#|[$time_local]|#|$msec|#|$real_host|#|$scheme|#|$x_qnm_cache|#|$request_method $real_request_uri $server_protocol|#|$status|#|$bytes_sent|#|$http_range|#|$sent_http_content_range|#|$http_referer|#|$http_user_agent|#|$http_x_forwarded_for|#|$request_time|#|$upstream_response_time|#|$upstream_connect_time|#|$upstream_header_time|#|$upstream_http_via|#|$upstream_addr|#|$upstream_http_content_type|#|$upstream_http_x_m_log|#|$master_reqid-$child_reqid|#|$real_hostname|#|$real_server_addr|#|$server_port|#|$first_byte_time|#|$sent_http_content_length|#|$request_length|#|$body_bytes_sent|#|$upstream_bytes_received|#|$upstream_status|#|$cache_level|#|$tcpinfo_rtt|#|$tcpinfo_rttvar|#|$tcpinfo_snd_cwnd|#|$tcpinfo_rcv_space|#|$connection_requests|#|$category|#|$internal_code|#|$http_x_from_cdn|#|$http_x_from_fsrcproxy|#|$response_from|#|$access_time|#|$access_end_time|#|$first_body_bytes_time|#|$log_time|#|$idc_name|#|remote_from=$remote_from|#|Qypid=$qypid|#|Qyid=$qyid|#|request_completion=$request_completion|#|c_auth_time=$c_auth_time|#|$http_content_type|#|$remote_port|#|redirect_to=$redirect_to|#|$cache_status|#|$src_ip|#|$device_level|#|-|#|$remote_user|#|$sent_http_last_modified|#|$nat_server_addr|#|options=$options"
		p := NewSpliteParser(pattern, "|#|")
		dict, _ := p.ParseString(input)
		convey.ShouldEqual(len(dict.Fields()), 67)
	})
}
