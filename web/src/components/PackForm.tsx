import { useMemo, useState } from "react";

type Props = {
  availableSizes?: number[];
  onCalculate: (quantity: number, packsOverride?: number[]) => Promise<void> | void;
  disabled?: boolean;
};

export default function PackForm({ availableSizes = [], onCalculate, disabled }: Props) {
  const [quantity, setQuantity] = useState<string>("");
  const [override, setOverride] = useState<string>(""); // ex: "250,500,1000"
  const [error, setError] = useState<string | null>(null);

  const parsedOverride = useMemo(() => {
    if (!override.trim()) return undefined;
    const uniq = new Set<number>();
    for (const tok of override.split(/[,\s;]+/)) {
      if (!tok.trim()) continue;
      const n = Number(tok.trim());
      if (!Number.isInteger(n) || n <= 0) return null; // inválido
      uniq.add(n);
    }
    return Array.from(uniq).sort((a, b) => a - b);
  }, [override]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    const q = Number(quantity);
    if (!Number.isInteger(q) || q <= 0) {
      setError("Informe uma quantidade inteira maior que 0.");
      return;
    }
    if (parsedOverride === null) {
      setError("Packs override inválidos. Use inteiros positivos, separados por vírgula.");
      return;
    }

    await onCalculate(q, parsedOverride ?? undefined);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <label className="block text-sm font-medium mb-1">Quantidade</label>
        <input
          inputMode="numeric"
          pattern="[0-9]*"
          value={quantity}
          onChange={(e) => setQuantity(e.target.value)}
          placeholder="Ex.: 12001"
          className="w-full rounded-md border px-3 py-2 outline-none focus:ring-2"
        />
      </div>

      <div>
        <div className="flex items-center justify-between">
          <label className="block text-sm font-medium mb-1">Override (opcional)</label>
          {availableSizes.length > 0 && (
            <span className="text-xs text-gray-500">
              Padrão: {availableSizes.join(", ")}
            </span>
          )}
        </div>
        <input
          value={override}
          onChange={(e) => setOverride(e.target.value)}
          placeholder="Ex.: 250, 500, 1000"
          className="w-full rounded-md border px-3 py-2 outline-none focus:ring-2"
        />
        <p className="text-xs text-gray-500 mt-1">
          Inteiros positivos separados por vírgula. Se preencher, substitui os tamanhos configurados.
        </p>
      </div>

      {error && <p className="text-sm text-red-600">{error}</p>}

      <button
        type="submit"
        disabled={disabled}
        className="rounded-md bg-black text-white px-4 py-2 disabled:opacity-50"
      >
        Calcular
      </button>
    </form>
  );
}