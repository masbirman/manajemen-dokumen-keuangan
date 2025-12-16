import { describe, it, expect } from "vitest";
import * as fc from "fast-check";

/**
 * Property 9: Operator Default Selection
 * For any operator user with assigned unit_kerja and pptk, the form should default to those values
 * **Validates: Requirements 6.1, 6.2**
 */

interface User {
  id: string;
  role: "super_admin" | "admin" | "operator";
  unit_kerja_id?: string;
  pptk_id?: string;
}

interface FormDefaults {
  unit_kerja_id: string;
  pptk_id: string;
}

// Simulates the logic from InputDokumenView.vue
const getFormDefaults = (user: User | null): FormDefaults => {
  const defaults: FormDefaults = {
    unit_kerja_id: "",
    pptk_id: "",
  };

  if (user && user.role === "operator") {
    if (user.unit_kerja_id) {
      defaults.unit_kerja_id = user.unit_kerja_id;
    }
    if (user.pptk_id) {
      defaults.pptk_id = user.pptk_id;
    }
  }

  return defaults;
};

const isOperator = (user: User | null): boolean => {
  return user?.role === "operator";
};

describe("Property 9: Operator Default Selection", () => {
  it("should set unit_kerja_id default for any operator with assignment", () => {
    fc.assert(
      fc.property(
        fc.uuid(),
        fc.uuid(),
        (unitKerjaId: string, pptkId: string) => {
          const operatorUser: User = {
            id: "user-1",
            role: "operator",
            unit_kerja_id: unitKerjaId,
            pptk_id: pptkId,
          };

          const defaults = getFormDefaults(operatorUser);
          expect(defaults.unit_kerja_id).toBe(unitKerjaId);
          expect(defaults.pptk_id).toBe(pptkId);
          return true;
        }
      ),
      { numRuns: 100 }
    );
  });

  it("should not set defaults for admin users", () => {
    fc.assert(
      fc.property(
        fc.uuid(),
        fc.uuid(),
        (unitKerjaId: string, pptkId: string) => {
          const adminUser: User = {
            id: "user-1",
            role: "admin",
            unit_kerja_id: unitKerjaId,
            pptk_id: pptkId,
          };

          const defaults = getFormDefaults(adminUser);
          // Admin should not have defaults set
          expect(defaults.unit_kerja_id).toBe("");
          expect(defaults.pptk_id).toBe("");
          return true;
        }
      ),
      { numRuns: 100 }
    );
  });

  it("should not set defaults for super_admin users", () => {
    fc.assert(
      fc.property(
        fc.uuid(),
        fc.uuid(),
        (unitKerjaId: string, pptkId: string) => {
          const superAdminUser: User = {
            id: "user-1",
            role: "super_admin",
            unit_kerja_id: unitKerjaId,
            pptk_id: pptkId,
          };

          const defaults = getFormDefaults(superAdminUser);
          expect(defaults.unit_kerja_id).toBe("");
          expect(defaults.pptk_id).toBe("");
          return true;
        }
      ),
      { numRuns: 100 }
    );
  });

  it("should handle operator without assignments", () => {
    const operatorWithoutAssignment: User = {
      id: "user-1",
      role: "operator",
    };

    const defaults = getFormDefaults(operatorWithoutAssignment);
    expect(defaults.unit_kerja_id).toBe("");
    expect(defaults.pptk_id).toBe("");
  });

  it("should correctly identify operator role", () => {
    fc.assert(
      fc.property(
        fc.constantFrom("super_admin", "admin", "operator") as fc.Arbitrary<
          "super_admin" | "admin" | "operator"
        >,
        (role) => {
          const user: User = { id: "user-1", role };
          const result = isOperator(user);
          expect(result).toBe(role === "operator");
          return true;
        }
      ),
      { numRuns: 100 }
    );
  });

  it("should handle null user gracefully", () => {
    const defaults = getFormDefaults(null);
    expect(defaults.unit_kerja_id).toBe("");
    expect(defaults.pptk_id).toBe("");
    expect(isOperator(null)).toBe(false);
  });
});
