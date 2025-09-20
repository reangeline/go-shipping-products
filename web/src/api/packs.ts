// src/api/packs.ts

export interface PackSizeResponse {
  sizes: number[];
}

export interface CalculateResponse {
  itemsByPack: Record<number, number>;
  totalItems: number;
  totalPacks: number;
  leftover: number;
}

const BASE_URL = "/v1";

/**
 * Busca a lista de pack sizes configurados no servidor
 */
export async function getPackSizes(): Promise<PackSizeResponse> {
  const res = await fetch(`${BASE_URL}/packsizes`);
  if (!res.ok) throw new Error(`Erro ao buscar packs: ${res.statusText}`);
  return res.json();
}

/**
 * Calcula a melhor combinação de packs para uma quantidade
 */
export async function calculatePacks(quantity: number): Promise<CalculateResponse> {
  const res = await fetch(`${BASE_URL}/calculate`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ quantity }),
  });
  if (!res.ok) {
    const msg = await res.text();
    throw new Error(msg || `Erro no cálculo (HTTP ${res.status})`);
  }
  return res.json();
}