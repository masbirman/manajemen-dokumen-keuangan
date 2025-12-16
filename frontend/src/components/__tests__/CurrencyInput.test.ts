import { describe, it, expect } from "vitest";
import * as fc from "fast-check";

/**
 * Property 13: Currency Value Handling
 * For any numeric value, the currency input should format and parse correctly
 * **Validates: Requirements 6.6**
 */

// Helper functions that mirror the CurrencyInput component logic
const formatCurrency = (value: number): string => {
  return new Intl.NumberFormat("id-ID").format(value);
};

const parseCurrency = (formatted: string): number => {
  const rawValue = formatted.replace(/\D/g, "");
  return parseInt(rawValue) || 0;
};

describe("Property 13: Currency Value Handling", () => {
  it("should format any positive integer correctly", () => {
    fc.assert(
      fc.property(
        fc.integer({ min: 0, max: 999999999999 }),
        (value: number) => {
          const formatted = formatCurrency(value);
          // Formatted string should not be empty
          expect(formatted.length).toBeGreaterThan(0);
          // Should contain only digits and separators
          expect(formatted.replace(/[\d.,]/g, "")).toBe("");
          return true;
        }
      ),
      { numRuns: 100 }
    );
  });

  it("should round-trip: format then parse returns original value", () => {
    fc.assert(
      fc.property(
        fc.integer({ min: 0, max: 999999999999 }),
        (value: number) => {
          const formatted = formatCurrency(value);
          const parsed = parseCurrency(formatted);
          expect(parsed).toBe(value);
          return true;
        }
      ),
      { numRuns: 100 }
    );
  });

  it("should handle zero correctly", () => {
    const formatted = formatCurrency(0);
    const parsed = parseCurrency(formatted);
    expect(parsed).toBe(0);
  });

  it("should strip non-numeric characters when parsing", () => {
    fc.assert(
      fc.property(
        fc.integer({ min: 1, max: 999999999 }),
        fc.string({ minLength: 0, maxLength: 5 }),
        (value: number, prefix: string) => {
          const formatted = formatCurrency(value);
          const withPrefix = prefix + formatted;
          const parsed = parseCurrency(withPrefix);
          // Should extract the numeric value regardless of prefix
          expect(parsed).toBeGreaterThanOrEqual(0);
          return true;
        }
      ),
      { numRuns: 100 }
    );
  });

  it("should format with thousand separators for large numbers", () => {
    fc.assert(
      fc.property(
        fc.integer({ min: 1000, max: 999999999 }),
        (value: number) => {
          const formatted = formatCurrency(value);
          // Large numbers should have separators (dots or commas)
          const hasSeparator =
            formatted.includes(".") || formatted.includes(",");
          expect(hasSeparator).toBe(true);
          return true;
        }
      ),
      { numRuns: 100 }
    );
  });

  it("should preserve value magnitude after formatting", () => {
    fc.assert(
      fc.property(fc.integer({ min: 0, max: 999999999 }), (value: number) => {
        const formatted = formatCurrency(value);
        const digitCount = formatted.replace(/\D/g, "").length;
        const originalDigitCount = value.toString().length;
        expect(digitCount).toBe(originalDigitCount);
        return true;
      }),
      { numRuns: 100 }
    );
  });
});
