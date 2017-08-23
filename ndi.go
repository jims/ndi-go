/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package ndi

import (
	"errors"
	"reflect"
	"syscall"
	"unsafe"
)

var (
	ndiSharedLibrary syscall.Handle
	funcPtrs         *ndiLIBv3
)

func goString(p uintptr) string {
	var len int
	for n := p; *(*byte)(unsafe.Pointer(n)) != 0; n++ {
		len++
	}

	h := &reflect.SliceHeader{uintptr(unsafe.Pointer(p)), len, len + 1}
	return string(*(*[]byte)(unsafe.Pointer(h)))
}

func LoadAndInitialize(path string) error {
	var err error
	if ndiSharedLibrary, err = syscall.LoadLibrary(path); err != nil {
		return err
	}

	var ndiLoadProc uintptr
	if ndiLoadProc, err = syscall.GetProcAddress(ndiSharedLibrary, "NDIlib_v3_load"); err != nil {
		syscall.FreeLibrary(ndiSharedLibrary)
		return err
	}

	var (
		ret uintptr
		eno syscall.Errno
	)

	if ret, _, eno = syscall.Syscall(ndiLoadProc, 0, 0, 0, 0); eno != 0 {
		syscall.FreeLibrary(ndiSharedLibrary)
		return eno
	}

	funcPtrs = (*ndiLIBv3)(unsafe.Pointer(ret))
	if funcPtrs == nil {
		return errors.New("failed to load library procs")
	}

	if _, _, eno = syscall.Syscall(funcPtrs.NDIlibInitialize, 0, 0, 0, 0); eno != 0 {
		syscall.FreeLibrary(ndiSharedLibrary)
		return eno
	}

	if ret == 0 {
		return errors.New("unable to initialize library")
	}
	return nil
}

func UnloadAndDestroy() {
	if ndiSharedLibrary == 0 {
		return
	}
	if funcPtrs != nil {
		syscall.Syscall(funcPtrs.NDIlibDestroy, 0, 0, 0, 0)
	}
	syscall.FreeLibrary(ndiSharedLibrary)
}

func Version() string {
	ret, _, _ := syscall.Syscall(funcPtrs.NDIlibVersion, 0, 0, 0, 0)
	return goString(ret)
}

type ndiLIBv3 struct {
	// V1.5
	NDIlibInitialize, //bool(*NDIlib_initialize)(void)
	NDIlibDestroy, //void(*NDIlib_destroy)(void)
	NDIlibVersion, //const char* (*NDIlib_version)(void)
	NDIlibIsSupportedCPU, //bool(*NDIlib_is_supported_CPU)(void)
	NDIlibFindCreate, //PROCESSINGNDILIB_DEPRECATED NDIlib_find_instance_t(*NDIlib_find_create)(const NDIlib_find_create_t* p_create_settings)
	NDIlibFindCreateV2, //NDIlib_find_instance_t(*NDIlib_find_create_v2)(const NDIlib_find_create_t* p_create_settings)
	NDIlibFindDestroy, //void(*NDIlib_find_destroy)(NDIlib_find_instance_t p_instance)
	NDIlibFindGetSources, //const NDIlib_source_t* (*NDIlib_find_get_sources)(NDIlib_find_instance_t p_instance, uint32_t* p_no_sources, uint32_t timeout_in_ms)
	NDIlibSendCreate, //NDIlib_send_instance_t(*NDIlib_send_create)(const NDIlib_send_create_t* p_create_settings)
	NDIlibSendDestroy, //void(*NDIlib_send_destroy)(NDIlib_send_instance_t p_instance)
	NDIlibSendSendVideo, //PROCESSINGNDILIB_DEPRECATED void(*NDIlib_send_send_video)(NDIlib_send_instance_t p_instance, const NDIlib_video_frame_t* p_video_data)
	NDIlibSendSendVideoAsync, //PROCESSINGNDILIB_DEPRECATED void(*NDIlib_send_send_video_async)(NDIlib_send_instance_t p_instance, const NDIlib_video_frame_t* p_video_data)
	NDIlibSendSendAudio, //PROCESSINGNDILIB_DEPRECATED void(*NDIlib_send_send_audio)(NDIlib_send_instance_t p_instance, const NDIlib_audio_frame_t* p_audio_data)
	NDIlibSendSendMetadata, //void(*NDIlib_send_send_metadata)(NDIlib_send_instance_t p_instance, const NDIlib_metadata_frame_t* p_metadata)
	NDIlibSendCapture, //NDIlib_frame_type_e(*NDIlib_send_capture)(NDIlib_send_instance_t p_instance, NDIlib_metadata_frame_t* p_metadata, uint32_t timeout_in_ms)
	NDIlibSendFreeMetadata, //void(*NDIlib_send_free_metadata)(NDIlib_send_instance_t p_instance, const NDIlib_metadata_frame_t* p_metadata)
	NDIlibSendGetTally, //bool(*NDIlib_send_get_tally)(NDIlib_send_instance_t p_instance, NDIlib_tally_t* p_tally, uint32_t timeout_in_ms)
	NDIlibSendGetNoConnections, //int(*NDIlib_send_get_no_connections)(NDIlib_send_instance_t p_instance, uint32_t timeout_in_ms)
	NDIlibSendClearConnectionMetadata, //void(*NDIlib_send_clear_connection_metadata)(NDIlib_send_instance_t p_instance)
	NDIlibSendAddConnectionMetadata, //void(*NDIlib_send_add_connection_metadata)(NDIlib_send_instance_t p_instance, const NDIlib_metadata_frame_t* p_metadata)
	NDIlibSendSetFailover, //void(*NDIlib_send_set_failover)(NDIlib_send_instance_t p_instance, const NDIlib_source_t* p_failover_source)
	NDIlibRecvCreateV2, //NDIlib_recv_instance_t(*NDIlib_recv_create_v2)(const NDIlib_recv_create_t* p_create_settings)
	NDIlibRecvCreate, //NDIlib_recv_instance_t(*NDIlib_recv_create)(const NDIlib_recv_create_t* p_create_settings)
	NDIlibRecvDestroy, //void(*NDIlib_recv_destroy)(NDIlib_recv_instance_t p_instance)
	NDIlibRecvCapture, //PROCESSINGNDILIB_DEPRECATED NDIlib_frame_type_e(*NDIlib_recv_capture)(NDIlib_recv_instance_t p_instance, NDIlib_video_frame_t* p_video_data, NDIlib_audio_frame_t* p_audio_data, NDIlib_metadata_frame_t* p_metadata, uint32_t timeout_in_ms)
	NDIlibRecvFreeVideo, //PROCESSINGNDILIB_DEPRECATED void(*NDIlib_recv_free_video)(NDIlib_recv_instance_t p_instance, const NDIlib_video_frame_t* p_video_data)
	NDIlibRecvFreeAudio, //PROCESSINGNDILIB_DEPRECATED void(*NDIlib_recv_free_audio)(NDIlib_recv_instance_t p_instance, const NDIlib_audio_frame_t* p_audio_data)
	NDIlibRecvFreeMetadata, //void(*NDIlib_recv_free_metadata)(NDIlib_recv_instance_t p_instance, const NDIlib_metadata_frame_t* p_metadata)
	NDIlibRecvSendMetadata, //bool(*NDIlib_recv_send_metadata)(NDIlib_recv_instance_t p_instance, const NDIlib_metadata_frame_t* p_metadata)
	NDIlibRecvSetTally, //bool(*NDIlib_recv_set_tally)(NDIlib_recv_instance_t p_instance, const NDIlib_tally_t* p_tally)
	NDIlibRecvGetPerformance, //void(*NDIlib_recv_get_performance)(NDIlib_recv_instance_t p_instance, NDIlib_recv_performance_t* p_total, NDIlib_recv_performance_t* p_dropped)
	NDIlibRecvGetQueue, //void(*NDIlib_recv_get_queue)(NDIlib_recv_instance_t p_instance, NDIlib_recv_queue_t* p_total)
	NDIlibRecvClearConnectionMetadata, //void(*NDIlib_recv_clear_connection_metadata)(NDIlib_recv_instance_t p_instance)
	NDIlibRecvAddConnectionMetadata, //void(*NDIlib_recv_add_connection_metadata)(NDIlib_recv_instance_t p_instance, const NDIlib_metadata_frame_t* p_metadata)
	NDIlibRecvGetNoConnections, //int(*NDIlib_recv_get_no_connections)(NDIlib_recv_instance_t p_instance)
	NDIlibRoutingCreate, //NDIlib_routing_instance_t(*NDIlib_routing_create)(const NDIlib_routing_create_t* p_create_settings)
	NDIlibRoutingDestroy, //void(*NDIlib_routing_destroy)(NDIlib_routing_instance_t p_instance)
	NDIlibRoutingChange, //bool(*NDIlib_routing_change)(NDIlib_routing_instance_t p_instance, const NDIlib_source_t* p_source)
	NDIlibRoutingClear, //bool(*NDIlib_routing_clear)(NDIlib_routing_instance_t p_instance)
	NDIlibUtilSendSendAudioInterleaved16s, //void(*NDIlib_util_send_send_audio_interleaved_16s)(NDIlib_send_instance_t p_instance, const NDIlib_audio_frame_interleaved_16s_t* p_audio_data)
	NDIlibUtilAudioToInterleaved16s, //PROCESSINGNDILIB_DEPRECATED void(*NDIlib_util_audio_to_interleaved_16s)(const NDIlib_audio_frame_t* p_src, NDIlib_audio_frame_interleaved_16s_t* p_dst)
	NDIlibUtilAudioFromInterleaved16s, //PROCESSINGNDILIB_DEPRECATED void(*NDIlib_util_audio_from_interleaved_16s)(const NDIlib_audio_frame_interleaved_16s_t* p_src, NDIlib_audio_frame_t* p_dst)

	// V2
	NDIlibFindWaitForSources, //bool(*NDIlib_find_wait_for_sources)(NDIlib_find_instance_t p_instance, uint32_t timeout_in_ms)
	NDIlibFindGetCurrentSources, //const NDIlib_source_t* (*NDIlib_find_get_current_sources)(NDIlib_find_instance_t p_instance, uint32_t* p_no_sources)
	NDIlibUtilAudioToInterleaved32f, //PROCESSINGNDILIB_DEPRECATED void(*NDIlib_util_audio_to_interleaved_32f)(const NDIlib_audio_frame_t* p_src, NDIlib_audio_frame_interleaved_32f_t* p_dst)
	NDIlibUtilAudioFromInterleaved32f, //PROCESSINGNDILIB_DEPRECATED void(*NDIlib_util_audio_from_interleaved_32f)(const NDIlib_audio_frame_interleaved_32f_t* p_src, NDIlib_audio_frame_t* p_dst)
	NDIlibUtilSendSendAudioInterleaved32f, //void(*NDIlib_util_send_send_audio_interleaved_32f)(NDIlib_send_instance_t p_instance, const NDIlib_audio_frame_interleaved_32f_t* p_audio_data)

	// V3
	NDIlibRecvFreeVideoV2, //void(*NDIlib_recv_free_video_v2)(NDIlib_recv_instance_t p_instance, const NDIlib_video_frame_v2_t* p_video_data)
	NDIlibRecvFreeAudioV2, //void(*NDIlib_recv_free_audio_v2)(NDIlib_recv_instance_t p_instance, const NDIlib_audio_frame_v2_t* p_audio_data)
	NDIlibRecvCaptureV2, //NDIlib_frame_type_e(*NDIlib_recv_capture_v2)(NDIlib_recv_instance_t p_instance, NDIlib_video_frame_v2_t* p_video_data, NDIlib_audio_frame_v2_t* p_audio_data, NDIlib_metadata_frame_t* p_metadata, uint32_t timeout_in_ms)
	NDIlibSendSendVideoV2, //void(*NDIlib_send_send_video_v2)(NDIlib_send_instance_t p_instance, const NDIlib_video_frame_v2_t* p_video_data)
	NDIlibSendSendVideoAsyncV2, //void(*NDIlib_send_send_video_async_v2)(NDIlib_send_instance_t p_instance, const NDIlib_video_frame_v2_t* p_video_data)
	NDIlibSendSendAudioV2, //void(*NDIlib_send_send_audio_v2)(NDIlib_send_instance_t p_instance, const NDIlib_audio_frame_v2_t* p_audio_data)
	NDIlibUtilAudioToInterleaved16sV2, //void(*NDIlib_util_audio_to_interleaved_16s_v2)(const NDIlib_audio_frame_v2_t* p_src, NDIlib_audio_frame_interleaved_16s_t* p_dst)
	NDIlibUtilAudioFromInterleaved16sV2, //void(*NDIlib_util_audio_from_interleaved_16s_v2)(const NDIlib_audio_frame_interleaved_16s_t* p_src, NDIlib_audio_frame_v2_t* p_dst)
	NDIlibUtilAudioToInterleaved32fV2, //void(*NDIlib_util_audio_to_interleaved_32f_v2)(const NDIlib_audio_frame_v2_t* p_src, NDIlib_audio_frame_interleaved_32f_t* p_dst)
	NDIlibUtilAudioFromInterleaved32fV2, //void(*NDIlib_util_audio_from_interleaved_32f_v2)(const NDIlib_audio_frame_interleaved_32f_t* p_src, NDIlib_audio_frame_v2_t* p_dst)

	// V3.01
	NDIlibRecvFreeString, //void(*NDIlib_recv_free_string)(NDIlib_recv_instance_t p_instance, const char* p_string)
	NDIlibRecvPtzIsSupported, //bool(*NDIlib_recv_ptz_is_supported)(NDIlib_recv_instance_t p_instance)
	NDIlibRecvRecordingIsSupported, //bool(*NDIlib_recv_recording_is_supported)(NDIlib_recv_instance_t p_instance)
	NDIlibRecvGetWebControl, //const char*(*NDIlib_recv_get_web_control)(NDIlib_recv_instance_t p_instance)
	NDIlibRecvPtzZoom, //bool(*NDIlib_recv_ptz_zoom)(NDIlib_recv_instance_t p_instance, const float zoom_value)
	NDIlibRecvPtzZoomSpeed, //bool(*NDIlib_recv_ptz_zoom_speed)(NDIlib_recv_instance_t p_instance, const float zoom_speed)
	NDIlibRecvPtzPanTilt, //bool(*NDIlib_recv_ptz_pan_tilt)(NDIlib_recv_instance_t p_instance, const float pan_value, const float tilt_value)
	NDIlibRecvPtzPanTiltSpeed, //bool(*NDIlib_recv_ptz_pan_tilt_speed)(NDIlib_recv_instance_t p_instance, const float pan_speed, const float tilt_speed)
	NDIlibRecvPtzStorePreset, //bool(*NDIlib_recv_ptz_store_preset)(NDIlib_recv_instance_t p_instance, const int preset_no)
	NDIlibRecvPtzRecallPreset, //bool(*NDIlib_recv_ptz_recall_preset)(NDIlib_recv_instance_t p_instance, const int preset_no, const float speed)
	NDIlibRecvPtzAutoFocus, //bool(*NDIlib_recv_ptz_auto_focus)(NDIlib_recv_instance_t p_instance)
	NDIlibRecvPtzFocus, //bool(*NDIlib_recv_ptz_focus)(NDIlib_recv_instance_t p_instance, const float focus_value)
	NDIlibRecvPtzFocusSpeed, //bool(*NDIlib_recv_ptz_focus_speed)(NDIlib_recv_instance_t p_instance, const float focus_speed)
	NDIlibRecvPtzWhiteBalanceAuto, //bool(*NDIlib_recv_ptz_white_balance_auto)(NDIlib_recv_instance_t p_instance)
	NDIlibRecvPtzWhiteBalanceIndoor, //bool(*NDIlib_recv_ptz_white_balance_indoor)(NDIlib_recv_instance_t p_instance)
	NDIlibRecvPtzWhiteBalanceOutdoor, //bool(*NDIlib_recv_ptz_white_balance_outdoor)(NDIlib_recv_instance_t p_instance)
	NDIlibRecvPtzWhiteBalanceOneshot, //bool(*NDIlib_recv_ptz_white_balance_oneshot)(NDIlib_recv_instance_t p_instance)
	NDIlibRecvPtzWhiteBalanceManual, //bool(*NDIlib_recv_ptz_white_balance_manual)(NDIlib_recv_instance_t p_instance, const float red, const float blue)
	NDIlibRecvPtzExposureAuto, //bool(*NDIlib_recv_ptz_exposure_auto)(NDIlib_recv_instance_t p_instance)
	NDIlibRecvPtzExposureManual, //bool(*NDIlib_recv_ptz_exposure_manual)(NDIlib_recv_instance_t p_instance, const float exposure_level)
	NDIlibRecvRecordingStart, //bool(*NDIlib_recv_recording_start)(NDIlib_recv_instance_t p_instance, const char* p_filename_hint)
	NDIlibRecvRecordingStop, //bool(*NDIlib_recv_recording_stop)(NDIlib_recv_instance_t p_instance)
	NDIlibRecvRecordingSetAudioLevel, //bool(*NDIlib_recv_recording_set_audio_level)(NDIlib_recv_instance_t p_instance, const float level_dB)
	NDIlibRecvRecordingIsRecording, //bool(*NDIlib_recv_recording_is_recording)(NDIlib_recv_instance_t p_instance)
	NDIlibRecvRecordingGetFilename, //const char*(*NDIlib_recv_recording_get_filename)(NDIlib_recv_instance_t p_instance)
	NDIlibRecvRecordingGetError, //const char*(*NDIlib_recv_recording_get_error)(NDIlib_recv_instance_t p_instance)
	NDIlibRecvRecordingGetTimes uintptr //bool(*NDIlib_recv_recording_get_times)(NDIlib_recv_instance_t p_instance, NDIlib_recv_recording_time_t* p_times)
}
