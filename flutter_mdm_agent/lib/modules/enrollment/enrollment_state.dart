/// Represents the current enrollment status of the device.
enum EnrollmentStatus {
  /// Device has not been enrolled yet.
  notEnrolled,

  /// Enrollment is in progress.
  enrolling,

  /// Device is successfully enrolled.
  enrolled,

  /// Enrollment failed with an error.
  error,
}
