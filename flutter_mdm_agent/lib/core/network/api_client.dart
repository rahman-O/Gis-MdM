import 'package:dio/dio.dart';
import '../utils/logger.dart';

/// HTTP client for communicating with the MDM server.
/// Supports retry, logging, and certificate pinning (Phase 2).
class ApiClient {
  late final Dio _dio;
  String? _baseUrl;

  ApiClient() {
    _dio = Dio(BaseOptions(
      connectTimeout: const Duration(seconds: 15),
      receiveTimeout: const Duration(seconds: 30),
      sendTimeout: const Duration(seconds: 15),
    ));
    _dio.interceptors.add(_LoggingInterceptor());
    _dio.interceptors.add(_RetryInterceptor(_dio));
  }

  void configure({required String baseUrl}) {
    _baseUrl = baseUrl.replaceAll(RegExp(r'/+$'), '');
    _dio.options.baseUrl = _baseUrl!;
    Logger.info('API client configured: $_baseUrl', 'ApiClient');
  }

  bool get isConfigured => _baseUrl != null && _baseUrl!.isNotEmpty;

  Future<Response<T>> get<T>(String path, {Map<String, dynamic>? queryParameters, Map<String, String>? headers}) {
    return _dio.get<T>(path, queryParameters: queryParameters, options: Options(headers: headers));
  }

  Future<Response<T>> post<T>(String path, {dynamic data, Map<String, String>? headers}) {
    return _dio.post<T>(path, data: data, options: Options(headers: headers));
  }

  Future<Response<T>> put<T>(String path, {dynamic data, Map<String, String>? headers}) {
    return _dio.put<T>(path, data: data, options: Options(headers: headers));
  }

  Future<Response> download(String url, String savePath) {
    return _dio.download(url, savePath);
  }
}

class _LoggingInterceptor extends Interceptor {
  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) {
    Logger.debug('→ ${options.method} ${options.path}', 'HTTP');
    handler.next(options);
  }

  @override
  void onResponse(Response response, ResponseInterceptorHandler handler) {
    Logger.debug('← ${response.statusCode} ${response.requestOptions.path}', 'HTTP');
    handler.next(response);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    Logger.error('✗ ${err.requestOptions.method} ${err.requestOptions.path}: ${err.message}', err, null, 'HTTP');
    handler.next(err);
  }
}

class _RetryInterceptor extends Interceptor {
  final Dio _dio;
  static const int _maxRetries = 3;

  _RetryInterceptor(this._dio);

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    final retryCount = err.requestOptions.extra['retryCount'] as int? ?? 0;
    if (retryCount >= _maxRetries) {
      handler.next(err);
      return;
    }
    if (_shouldRetry(err)) {
      await Future.delayed(Duration(milliseconds: 1000 * (retryCount + 1)));
      err.requestOptions.extra['retryCount'] = retryCount + 1;
      try {
        final response = await _dio.fetch(err.requestOptions);
        handler.resolve(response);
      } catch (e) {
        handler.next(err);
      }
    } else {
      handler.next(err);
    }
  }

  bool _shouldRetry(DioException err) {
    return err.type == DioExceptionType.connectionTimeout ||
        err.type == DioExceptionType.receiveTimeout ||
        err.type == DioExceptionType.connectionError ||
        (err.response?.statusCode != null && err.response!.statusCode! >= 500);
  }
}
