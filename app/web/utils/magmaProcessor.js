import { calculateXOR, paarAlgorithm, bpAlgorithm, bpAlgorithmWithDepth } from './algorithms';


  export function processMagmaFile(content) {
  const results = [];
  const matrixRegex = /\[([0-1 ]+)\]/g;
  const matrices = content.split('-----------------');

  matrices.forEach((matrixBlock, index) => {
    if (!matrixBlock.trim()) return;

    const binaryMatrix = [];
    let match;
    while ((match = matrixRegex.exec(matrixBlock)) !== null) {
      binaryMatrix.push(match[1].trim());
    }

    if (binaryMatrix.length > 0) {
      const matrixString = binaryMatrix.join('\n');
      const xorCount = calculateXOR(matrixString);
      const paarResult = paarAlgorithm(matrixString);
      const bpResult = bpAlgorithm(matrixString);
      const bpDepthResult = bpAlgorithmWithDepth(matrixString, 2);

      results.push({
        matrixNumber: index + 1,
        binaryMatrix,
        xorCount,
        paarResult,
        bpResult,
        bpDepthResult
      });
    }
  });

  return results;
}

export function processMatrix(matrixText) {
  // Matrisi satırlara böl
  const rows = matrixText.trim().split('\n').map(row => row.trim()).filter(row => row);
  
  // Her satırı binary matris formatına dönüştür
  const binaryMatrix = rows.map(row => {
    // Köşeli parantezleri ve boşlukları temizle
    return row.replace(/[\[\]]/g, '').trim();
  });

  // Matrisi string formatına dönüştür
  const matrixString = binaryMatrix.join('\n');

  // Sonuçları hesapla
  const xorCount = calculateXOR(matrixString);
  const paarResult = paarAlgorithm(matrixString);
  const bpResult = bpAlgorithm(matrixString);
  const bpDepthResult = bpAlgorithmWithDepth(matrixString, 2);

  return {
    binaryMatrix,
    xorCount,
    paarResult,
    bpResult,
    bpDepthResult
  };
} 