import type { ProfileTreeAssignment } from '@/features/profiles/profileRolloutService'

/** True when another assignment's path is a strict prefix of this row's path. */
export function assignmentHasParentOverlap(
  row: ProfileTreeAssignment,
  all: ProfileTreeAssignment[]
): boolean {
  const path = row.treePath?.trim()
  if (!path) return false
  return all.some((other) => {
    if (other.assignmentId === row.assignmentId) return false
    const otherPath = other.treePath?.trim()
    if (!otherPath || otherPath === path) return false
    return path.startsWith(otherPath) && path.length > otherPath.length
  })
}

/** Parent + child both assigned to the same profile (nearest node wins on devices). */
export function hasParentChildAssignmentOverlap(assignments: ProfileTreeAssignment[]): boolean {
  return assignments.some((row) => assignmentHasParentOverlap(row, assignments))
}
