import { useEffect, useState } from "react";
import { getPackSizes, calculatePacks,  } from "../api/packs";

import type {CalculateResponse} from "../api/packs"

export default function PacksCalculator() {
  const [sizes, setSizes] = useState<number[]>([]);
  const [loadingSizes, setLoadingSizes] = useState(false);
  const [sizesError, setSizesError] = useState<string | null>(null);

  const [qty, setQty] = useState("12001");
  const [calcLoading, setCalcLoading] = useState(false);
  const [calcError, setCalcError] = useState<string | null>(null);
  const [result, setResult] = useState<CalculateResponse | null>(null);

  useEffect(() => {
    (async () => {
      try {
        setLoadingSizes(true);
        const s = await getPackSizes();
        setSizes(s.sizes);
      } catch (e: any) {
        setSizesError(e.message || String(e));
      } finally {
        setLoadingSizes(false);
      }
    })();
  }, []);

    const handleCalculate = async () => {
    setCalcError(null);
    setCalcLoading(true);
    setResult(null);

    const q = Number(qty);
    if (!Number.isInteger(q) || q <= 0) {
      setCalcLoading(false);
      setCalcError("Quantity must be a positive integer");
      return;
    }

    try {
      const resp = await calculatePacks(q);
      setResult(resp);
    } catch (e: any) {
      setCalcError(e?.message ?? "Unexpected error");
    } finally {
      setCalcLoading(false);
    }
  };


  return (
    <div className="min-h-screen  bg-white text-gray-900">
      <div className="max-w-3xl mx-auto p-6 space-y-8">
        <h1 className="text-3xl font-extrabold">Order Packs Calculator</h1>

        {/* Pack Sizes (read-only table) */}
        <section className="space-y-3">
          <h2 className="text-xl font-semibold">Pack Sizes</h2>
          <div className="rounded-lg border overflow-hidden">
            {loadingSizes ? (
              <div className="p-4 text-sm text-gray-500">Loading…</div>
            ) : sizesError ? (
              <div className="p-4 text-sm text-red-600">{sizesError}</div>
            ) : (
              <table className="w-full border-collapse">
                <thead>
                  <tr className="bg-gray-100 text-left">
                    <th className="border px-3 py-2">Pack</th>
                  </tr>
                </thead>
                <tbody>
                  {sizes.length === 0 ? (
                    <tr>
                      <td className="border px-3 py-2 text-sm text-gray-500">—</td>
                    </tr>
                  ) : (
                    sizes.map((size) => (
                      <tr key={size}>
                        <td className="border px-3 py-2">{size}</td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            )}
          </div>
        </section>

        {/* Calculate */}
        <section className="space-y-3">
          <h2 className="text-xl font-semibold">Calculate packs for order</h2>
          <div className="flex items-center gap-3">
            <label className="text-sm">Items:</label>
            <input
              className="w-48 rounded-md border px-3 py-2"
              inputMode="numeric"
              style={{backgroundColor: "#fff"}}
              pattern="[0-9]*"
              value={qty}
              onChange={(e) => setQty(e.target.value)}
            />
            <button
              type="button"
              onClick={handleCalculate}
              className="rounded-md bg-green-600 text-white px-4 py-2"
              disabled={calcLoading}
            >
              {calcLoading ? "Calculating…" : "Calculate"}
            </button>
          </div>

          {/* Results */}
          {calcError && <p className="text-sm text-red-600">{calcError}</p>}
          {result && (
            <div className="mt-3">
              <table className="w-full border-collapse">
                <thead>
                  <tr className="bg-gray-100 text-left">
                    <th className="border px-3 py-2">Pack</th>
                    <th className="border px-3 py-2">Quantity</th>
                  </tr>
                </thead>
                <tbody>
                  {Object.entries(result.itemsByPack)
                    .map(([size, count]) => [Number(size), count as number] as const)
                    .sort((a, b) => b[0] - a[0])
                    .map(([size, count]) => (
                      <tr key={size}>
                        <td className="border px-3 py-2">{size}</td>
                        <td className="border px-3 py-2">{count}</td>
                      </tr>
                    ))}
                </tbody>
              </table>
              <div className="mt-3 grid grid-cols-3 gap-3 text-sm">
                <div className="rounded-md border p-3">
                  <div className="text-gray-600">Total Items</div>
                  <div className="font-semibold">{result.totalItems}</div>
                </div>
                <div className="rounded-md border p-3">
                  <div className="text-gray-600">Total Packs</div>
                  <div className="font-semibold">{result.totalPacks}</div>
                </div>
                <div className="rounded-md border p-3">
                  <div className="text-gray-600">Leftover</div>
                  <div className="font-semibold">{result.leftover}</div>
                </div>
              </div>
            </div>
          )}
        </section>
      </div>
    </div>
  );
}