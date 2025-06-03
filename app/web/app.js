// Global variables
let currentPage = 1;
let currentLimit = 10;
let currentFilter = '';
let currentGroupFilter = '';
let currentMatrixId = null;
let currentFilters = {
    hamXorMin: null,
    hamXorMax: null,
    boyarXorMin: null,
    boyarXorMax: null,
    paarXorMin: null,
    paarXorMax: null,
    slpXorMin: null,
    slpXorMax: null
};

// Inverse Pairs functionality
let currentPairPage = 1;
let currentPairLimit = 25;
let currentPairGroupFilter = '';
let currentMaxCombinedXor = null;
let currentSortOrder = 'combined_asc';

// Initialize the application
document.addEventListener('DOMContentLoaded', function() {
    // Load URL parameters
    loadFromURL();
    loadMatrices();
    setupEventListeners();
    
    // Handle browser back/forward
    window.addEventListener('popstate', function(event) {
        loadFromURL();
        loadMatrices();
    });
});

// Load parameters from URL
function loadFromURL() {
    const urlParams = new URLSearchParams(window.location.search);
    
    // Load limit
    const limit = urlParams.get('limit');
    if (limit && !isNaN(limit) && limit > 0) {
        currentLimit = parseInt(limit);
        document.getElementById('pageSizeSelect').value = limit;
    }
    
    // Load search filter
    const search = urlParams.get('search');
    if (search) {
        currentFilter = search;
        document.getElementById('searchInput').value = search;
    }
    
    // Load group filter
    const group = urlParams.get('group');
    if (group) {
        currentGroupFilter = group;
        document.getElementById('groupSearchInput').value = group;
    }
    
    // Load range filters
    const hamXorMin = urlParams.get('hamXorMin');
    const hamXorMax = urlParams.get('hamXorMax');
    const boyarXorMin = urlParams.get('boyarXorMin');
    const boyarXorMax = urlParams.get('boyarXorMax');
    const paarXorMin = urlParams.get('paarXorMin');
    const paarXorMax = urlParams.get('paarXorMax');
    const slpXorMin = urlParams.get('slpXorMin');
    const slpXorMax = urlParams.get('slpXorMax');
    
    if (hamXorMin) {
        currentFilters.hamXorMin = parseInt(hamXorMin);
        document.getElementById('hamXorMin').value = hamXorMin;
    }
    if (hamXorMax) {
        currentFilters.hamXorMax = parseInt(hamXorMax);
        document.getElementById('hamXorMax').value = hamXorMax;
    }
    if (boyarXorMin) {
        currentFilters.boyarXorMin = parseInt(boyarXorMin);
        document.getElementById('boyarXorMin').value = boyarXorMin;
    }
    if (boyarXorMax) {
        currentFilters.boyarXorMax = parseInt(boyarXorMax);
        document.getElementById('boyarXorMax').value = boyarXorMax;
    }
    if (paarXorMin) {
        currentFilters.paarXorMin = parseInt(paarXorMin);
        document.getElementById('paarXorMin').value = paarXorMin;
    }
    if (paarXorMax) {
        currentFilters.paarXorMax = parseInt(paarXorMax);
        document.getElementById('paarXorMax').value = paarXorMax;
    }
    if (slpXorMin) {
        currentFilters.slpXorMin = parseInt(slpXorMin);
        document.getElementById('slpXorMin').value = slpXorMin;
    }
    if (slpXorMax) {
        currentFilters.slpXorMax = parseInt(slpXorMax);
        document.getElementById('slpXorMax').value = slpXorMax;
    }
}

// Update URL with current parameters
function updateURL() {
    const params = new URLSearchParams();
    
    // Add limit if not default
    if (currentLimit !== 10) {
        params.set('limit', currentLimit);
    }
    
    // Add search filter
    if (currentFilter) {
        params.set('search', currentFilter);
    }
    
    // Add group filter
    if (currentGroupFilter) {
        params.set('group', currentGroupFilter);
    }
    
    // Add range filters
    if (currentFilters.hamXorMin !== null) {
        params.set('hamXorMin', currentFilters.hamXorMin);
    }
    if (currentFilters.hamXorMax !== null) {
        params.set('hamXorMax', currentFilters.hamXorMax);
    }
    if (currentFilters.boyarXorMin !== null) {
        params.set('boyarXorMin', currentFilters.boyarXorMin);
    }
    if (currentFilters.boyarXorMax !== null) {
        params.set('boyarXorMax', currentFilters.boyarXorMax);
    }
    if (currentFilters.paarXorMin !== null) {
        params.set('paarXorMin', currentFilters.paarXorMin);
    }
    if (currentFilters.paarXorMax !== null) {
        params.set('paarXorMax', currentFilters.paarXorMax);
    }
    if (currentFilters.slpXorMin !== null) {
        params.set('slpXorMin', currentFilters.slpXorMin);
    }
    if (currentFilters.slpXorMax !== null) {
        params.set('slpXorMax', currentFilters.slpXorMax);
    }
    
    const newURL = window.location.pathname + (params.toString() ? '?' + params.toString() : '');
    window.history.pushState({}, '', newURL);
}

// Setup event listeners
function setupEventListeners() {
    // Search input
    const searchInput = document.getElementById('searchInput');
    let searchTimeout;
    searchInput.addEventListener('input', function() {
        clearTimeout(searchTimeout);
        searchTimeout = setTimeout(() => {
            currentFilter = this.value;
            currentPage = 1; // Reset to first page when searching
            loadMatrices();
        }, 500);
    });

    // Group search input
    const groupSearchInput = document.getElementById('groupSearchInput');
    let groupSearchTimeout;
    groupSearchInput.addEventListener('input', function() {
        clearTimeout(groupSearchTimeout);
        groupSearchTimeout = setTimeout(() => {
            currentGroupFilter = this.value;
            currentPage = 1; // Reset to first page when searching
            loadMatrices();
        }, 500);
    });

    // Add matrix form
    document.getElementById('addMatrixForm').addEventListener('submit', function(e) {
        e.preventDefault();
        addMatrix();
    });

    // Bulk upload form
    document.getElementById('bulkUploadForm').addEventListener('submit', function(e) {
        e.preventDefault();
        bulkUploadMatrices();
    });

    // Input method change
    document.querySelectorAll('input[name="inputMethod"]').forEach(radio => {
        radio.addEventListener('change', function() {
            toggleInputMethod(this.value);
        });
    });

    // Recalculate button
    document.getElementById('recalculateBtn').addEventListener('click', function() {
        recalculateMatrix();
    });

    // Calculate inverse button
    document.getElementById('calculateInverseBtn').addEventListener('click', function() {
        console.log('Calculate inverse button clicked!');
        calculateMatrixInverse();
    });

    // Page size selector
    document.getElementById('pageSizeSelect').addEventListener('change', function() {
        currentLimit = parseInt(this.value);
        currentPage = 1; // Reset to first page when changing page size
        loadMatrices();
    });

    // Pagination event delegation
    document.addEventListener('click', function(e) {
        if (e.target.matches('.page-link[data-page]')) {
            e.preventDefault();
            const page = parseInt(e.target.getAttribute('data-page'));
            if (page && page > 0) {
                changePage(page);
            }
        }
    });

    // Bulk Inverse Form
    const bulkInverseForm = document.getElementById('bulkInverseForm');
    if (bulkInverseForm) {
        bulkInverseForm.addEventListener('submit', function(e) {
            e.preventDefault();
            startBulkInverse();
        });
    }
    
    // Load stats when tab is shown
    const bulkInverseTab = document.getElementById('bulk-inverse-tab');
    if (bulkInverseTab) {
        bulkInverseTab.addEventListener('shown.bs.tab', function() {
            loadBulkInverseStats();
        });
    }

    // Tab change event listeners
    document.querySelectorAll('a[data-bs-toggle="tab"]').forEach(tab => {
        tab.addEventListener('shown.bs.tab', function (e) {
            const targetId = e.target.getAttribute('href').substring(1);
            
            if (targetId === 'matrices') {
                loadMatrices();
            } else if (targetId === 'bulk-inverse') {
                loadBulkInverseStats();
            } else if (targetId === 'inverse-pairs') {
                loadInversePairs();
            }
        });
    });
}

// Load matrices from API
async function loadMatrices() {
    try {
        const startTime = performance.now();
        showLoading('Matrisler yükleniyor...');
        
        const params = new URLSearchParams({
            page: currentPage,
            limit: currentLimit
        });
        
        if (currentFilter) {
            params.append('title', currentFilter);
        }

        if (currentGroupFilter) {
            params.append('group', currentGroupFilter);
        }

        // Add range filters
        if (currentFilters.hamXorMin !== null) {
            params.append('ham_xor_min', currentFilters.hamXorMin);
        }
        if (currentFilters.hamXorMax !== null) {
            params.append('ham_xor_max', currentFilters.hamXorMax);
        }
        if (currentFilters.boyarXorMin !== null) {
            params.append('boyar_xor_min', currentFilters.boyarXorMin);
        }
        if (currentFilters.boyarXorMax !== null) {
            params.append('boyar_xor_max', currentFilters.boyarXorMax);
        }
        if (currentFilters.paarXorMin !== null) {
            params.append('paar_xor_min', currentFilters.paarXorMin);
        }
        if (currentFilters.paarXorMax !== null) {
            params.append('paar_xor_max', currentFilters.paarXorMax);
        }
        if (currentFilters.slpXorMin !== null) {
            params.append('slp_xor_min', currentFilters.slpXorMin);
        }
        if (currentFilters.slpXorMax !== null) {
            params.append('slp_xor_max', currentFilters.slpXorMax);
        }

        // Add cache busting parameter every 30 seconds
        const cacheKey = Math.floor(Date.now() / 30000);
        params.append('_cache', cacheKey);

        const response = await fetch(`/api/matrices?${params}`, {
            method: 'GET',
            headers: {
                'Accept': 'application/json',
                'Cache-Control': 'max-age=30'
            }
        });
        
        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.message || 'Matrisler yüklenemedi');
        }

        const loadTime = performance.now() - startTime;
        console.log(`📊 API Response (${loadTime.toFixed(2)}ms):`, {
            total: data.total,
            page: data.page,
            total_pages: data.total_pages,
            currentPage: currentPage,
            matrices_count: data.matrices.length,
            load_time: `${loadTime.toFixed(2)}ms`
        });

        // If current page is greater than total pages, reset to page 1
        if (currentPage > data.total_pages && data.total_pages > 0) {
            console.log('Current page is greater than total pages, resetting to page 1');
            currentPage = 1;
            loadMatrices(); // Reload with page 1
            return;
        }

        // Sync currentPage with API response
        currentPage = data.page;

        displayMatrices(data.matrices);
        setupPagination(data.page, data.total_pages, data.total);
        
        // Update URL with current parameters
        updateURL();
        
        // Show performance info if load time is high
        if (loadTime > 1000) {
            console.warn(`⚠️ Slow API response: ${loadTime.toFixed(2)}ms`);
        }
        
    } catch (error) {
        console.error('Error loading matrices:', error);
        showAlert('Matrisler yüklenirken hata oluştu: ' + error.message, 'danger');
    } finally {
        hideLoading();
    }
}

// Display matrices in the UI
function displayMatrices(matrices) {
    const container = document.getElementById('matricesContainer');
    
    if (!matrices || matrices.length === 0) {
        container.innerHTML = `
            <div class="text-center py-5">
                <i class="fas fa-table fa-3x text-muted mb-3"></i>
                <h5 class="text-muted">Henüz matris bulunmuyor</h5>
                <p class="text-muted">Yeni matris eklemek için "Yeni Matris" sekmesini kullanın.</p>
            </div>
        `;
        return;
    }

    const html = matrices.map(matrix => `
        <div class="card mb-3">
            <div class="card-header d-flex justify-content-between align-items-center">
                <h6 class="mb-0">${escapeHtml(matrix.title)}</h6>
                <div>
                    <small class="text-muted me-3">ID: ${matrix.id}</small>
                    <button class="btn btn-sm btn-outline-primary" onclick="viewMatrix(${matrix.id})">
                        <i class="fas fa-eye me-1"></i>Detay
                    </button>
                </div>
            </div>
            <div class="card-body">
                <div class="row">
                    <div class="col-md-6">
                        <div class="matrix-display" style="max-height: 150px; overflow-y: auto;">
                            ${matrix.matrix_binary}
                        </div>
                    </div>
                    <div class="col-md-6">
                        <div class="row">
                            <div class="col-6">
                                <strong>Ham XOR:</strong> ${matrix.ham_xor_count}
                            </div>
                            <div class="col-6">
                                <strong>En Küçük XOR:</strong> ${matrix.smallest_xor || 'N/A'}
                            </div>
                        </div>
                        <div class="row mt-1">
                            <div class="col-6">
                                <strong>Grup:</strong> ${matrix.group || 'N/A'}
                            </div>
                            <div class="col-6">
                                <strong>Hex:</strong> <small>${matrix.matrix_hex}</small>
                            </div>
                        </div>
                        <div class="row mt-2">
                            <div class="col-4">
                                <div class="algorithm-result result-boyar">
                                    <strong>Boyar:</strong> ${matrix.boyar_xor_count || 'N/A'}
                                    ${matrix.boyar_depth ? `(D:${matrix.boyar_depth})` : ''}
                                </div>
                            </div>
                            <div class="col-4">
                                <div class="algorithm-result result-paar">
                                    <strong>Paar:</strong> ${matrix.paar_xor_count || 'N/A'}
                                </div>
                            </div>
                            <div class="col-4">
                                <div class="algorithm-result result-slp">
                                    <strong>SLP:</strong> ${matrix.slp_xor_count || 'N/A'}
                                </div>
                            </div>
                        </div>
                        <div class="row mt-2">
                            <div class="col-12">
                                <small class="text-muted">
                                    Oluşturulma: ${formatDate(matrix.created_at)}
                                    ${matrix.updated_at !== matrix.created_at ? `| Güncelleme: ${formatDate(matrix.updated_at)}` : ''}
                                </small>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    `).join('');

    container.innerHTML = html;
}

// Setup pagination
function setupPagination(currentPage, totalPages, totalItems) {
    const container = document.getElementById('paginationContainer');
    const pagination = document.getElementById('pagination');
    const paginationInfo = document.getElementById('paginationInfo');
    
    if (totalPages <= 1) {
        container.style.display = 'none';
        return;
    }
    
    container.style.display = 'block';
    
    // Update pagination info
    const startItem = (currentPage - 1) * currentLimit + 1;
    const endItem = Math.min(currentPage * currentLimit, totalItems);
    paginationInfo.textContent = `${startItem}-${endItem} / ${totalItems} matris gösteriliyor (Sayfa ${currentPage}/${totalPages})`;
    
    let html = '';
    
    // Previous button
    html += `
        <li class="page-item ${currentPage === 1 ? 'disabled' : ''}">
            <a class="page-link" href="#" ${currentPage > 1 ? `data-page="${currentPage - 1}"` : ''} aria-label="Previous">
                <span aria-hidden="true">&laquo;</span>
            </a>
        </li>
    `;
    
    // Page numbers
    const startPage = Math.max(1, currentPage - 2);
    const endPage = Math.min(totalPages, currentPage + 2);
    
    if (startPage > 1) {
        html += `<li class="page-item"><a class="page-link" href="#" data-page="1">1</a></li>`;
        if (startPage > 2) {
            html += `<li class="page-item disabled"><span class="page-link">...</span></li>`;
        }
    }
    
    for (let i = startPage; i <= endPage; i++) {
        html += `
            <li class="page-item ${i === currentPage ? 'active' : ''}">
                <a class="page-link" href="#" ${i !== currentPage ? `data-page="${i}"` : ''}>${i}</a>
            </li>
        `;
    }
    
    if (endPage < totalPages) {
        if (endPage < totalPages - 1) {
            html += `<li class="page-item disabled"><span class="page-link">...</span></li>`;
        }
        html += `<li class="page-item"><a class="page-link" href="#" data-page="${totalPages}">${totalPages}</a></li>`;
    }
    
    // Next button
    html += `
        <li class="page-item ${currentPage === totalPages ? 'disabled' : ''}">
            <a class="page-link" href="#" ${currentPage < totalPages ? `data-page="${currentPage + 1}"` : ''} aria-label="Next">
                <span aria-hidden="true">&raquo;</span>
            </a>
        </li>
    `;
    
    pagination.innerHTML = html;
}

// Change page
function changePage(page) {
    if (page < 1) return;
    currentPage = page;
    loadMatrices();
}

// Add new matrix
async function addMatrix() {
    try {
        const title = document.getElementById('matrixTitle').value;
        const group = document.getElementById('matrixGroup').value.trim();
        const matrixData = document.getElementById('matrixData').value.trim();
        const processImmediately = document.getElementById('processImmediately').checked;
        
        // Parse matrix data - support both simple format and JSON format
        let matrix;
        
        // Try simple format first, then JSON format
        try {
            // First try to parse as simple format: [0 1 0 1] per line
            matrix = parseSimpleMatrixFormat(matrixData);
        } catch (simpleError) {
            // If simple format fails, try JSON format
            try {
                matrix = JSON.parse(matrixData);
                if (!Array.isArray(matrix) || !Array.isArray(matrix[0])) {
                    throw new Error('JSON verisi 2D array formatında olmalı');
                }
            } catch (jsonError) {
                // Both formats failed, show more helpful error
                throw new Error(`Geçersiz matris formatı. Basit format hatası: ${simpleError.message}. JSON format hatası: ${jsonError.message}`);
            }
        }
        
        if (!Array.isArray(matrix) || !Array.isArray(matrix[0])) {
            throw new Error('Matris 2D array formatında olmalı');
        }
        
        showLoading('Matris kaydediliyor...');
        
        const endpoint = processImmediately ? '/api/matrices/process' : '/api/matrices';
        const response = await fetch(endpoint, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                title: title,
                group: group,
                matrix: matrix
            })
        });
        
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.message || 'Matris kaydedilemedi');
        }
        
        showAlert('Matris başarıyla kaydedildi!', 'success');
        
        // Clear form
        document.getElementById('addMatrixForm').reset();
        
        // Switch to matrices tab and reload
        const matricesTab = new bootstrap.Tab(document.getElementById('matrices-tab'));
        matricesTab.show();
        loadMatrices();
        
    } catch (error) {
        console.error('Error adding matrix:', error);
        showAlert('Matris eklenirken hata oluştu: ' + error.message, 'danger');
    } finally {
        hideLoading();
    }
}

// Parse simple matrix format: [0 1 0 1] per line
function parseSimpleMatrixFormat(matrixData) {
    const lines = matrixData.split('\n').map(line => line.trim()).filter(line => line.length > 0);
    
    if (lines.length === 0) {
        throw new Error('Boş matris verisi');
    }
    
    const matrix = [];
    
    for (let i = 0; i < lines.length; i++) {
        const line = lines[i];
        
        // Check if line starts with [ and ends with ]
        if (!line.startsWith('[') || !line.endsWith(']')) {
            throw new Error(`Satır ${i + 1}: Köşeli parantez ile başlayıp bitmelidir`);
        }
        
        // Extract content between brackets
        const content = line.slice(1, -1).trim();
        
        if (content.length === 0) {
            throw new Error(`Satır ${i + 1}: Boş satır`);
        }
        
        // Split by spaces and filter empty strings
        const elements = content.split(/\s+/).filter(el => el.length > 0);
        
        if (elements.length === 0) {
            throw new Error(`Satır ${i + 1}: Geçerli eleman bulunamadı`);
        }
        
        // Validate that all elements are 0 or 1
        for (let j = 0; j < elements.length; j++) {
            if (elements[j] !== '0' && elements[j] !== '1') {
                throw new Error(`Satır ${i + 1}, Eleman ${j + 1}: Sadece 0 ve 1 değerleri kabul edilir (${elements[j]} geçersiz)`);
            }
        }
        
        // Check that all rows have the same length
        if (matrix.length > 0 && elements.length !== matrix[0].length) {
            throw new Error(`Satır ${i + 1}: Tüm satırlar aynı uzunlukta olmalı (beklenen: ${matrix[0].length}, bulunan: ${elements.length})`);
        }
        
        matrix.push(elements);
    }
    
    return matrix;
}

// View matrix details
async function viewMatrix(id) {
    try {
        showLoading('Matris detayları yükleniyor...');
        
        const response = await fetch(`/api/matrices/${id}`);
        const matrix = await response.json();
        
        if (!response.ok) {
            throw new Error(matrix.message || 'Matris detayları alınamadı');
        }
        
        currentMatrixId = id;
        displayMatrixDetails(matrix);
        
        const modal = new bootstrap.Modal(document.getElementById('matrixModal'));
        modal.show();
        
    } catch (error) {
        console.error('Error viewing matrix:', error);
        showAlert('Matris detayları yüklenirken hata oluştu: ' + error.message, 'danger');
    } finally {
        hideLoading();
    }
}

// Display matrix details in modal
function displayMatrixDetails(matrix) {
    document.getElementById('matrixModalTitle').textContent = matrix.title;
    
    const modalBody = document.getElementById('matrixModalBody');
    
    const html = `
        <div class="row">
            <div class="col-md-6">
                <h6>Matris Verisi</h6>
                <div class="matrix-display" style="max-height: 300px; overflow-y: auto;">
                    ${matrix.matrix_binary}
                </div>
                <div class="mt-3">
                    <strong>Hex Gösterim:</strong><br>
                    <code style="word-break: break-all;">${matrix.matrix_hex}</code>
                </div>
                <div class="mt-3">
                    <strong>Hash:</strong><br>
                    <code>${matrix.matrix_hash}</code>
                </div>
                
                ${matrix.inverse_matrix_id ? `
                <div class="mt-3 p-3 border rounded bg-light">
                    <h6 class="text-primary">Ters Matris Bilgisi</h6>
                    <div class="mb-2">
                        <strong>Ters Matris ID:</strong> ${matrix.inverse_matrix_id}
                    </div>
                    <div class="mb-2">
                        <strong>Ters Matris Hash:</strong><br>
                        <code style="font-size: 0.8em;">${matrix.inverse_matrix_hash || 'N/A'}</code>
                    </div>
                    <button class="btn btn-sm btn-outline-primary" onclick="viewMatrix(${matrix.inverse_matrix_id})">
                        <i class="fas fa-eye"></i> Ters Matrisi Görüntüle
                    </button>
                </div>
                ` : ''}
            </div>
            <div class="col-md-6">
                <h6>Algoritma Sonuçları</h6>
                
                <div class="mb-3">
                    <strong>Ham XOR Sayısı:</strong> ${matrix.ham_xor_count}
                </div>
                
                <div class="mb-3">
                    <strong>En Küçük XOR:</strong> ${matrix.smallest_xor || 'Hesaplanmamış'}
                </div>
                
                <div class="mb-3">
                    <strong>Grup:</strong> ${matrix.group || 'Belirtilmemiş'}
                </div>
                
                <div class="algorithm-result result-boyar mb-3">
                    <strong>Boyar SLP:</strong><br>
                    XOR: ${matrix.boyar_xor_count || 'Hesaplanmamış'}<br>
                    Derinlik: ${matrix.boyar_depth || 'N/A'}<br>
                    ${matrix.boyar_program ? `<details><summary>Program</summary><pre>${JSON.stringify(JSON.parse(matrix.boyar_program), null, 2)}</pre></details>` : ''}
                </div>
                
                <div class="algorithm-result result-paar mb-3">
                    <strong>Paar Algoritması:</strong><br>
                    XOR: ${matrix.paar_xor_count || 'Hesaplanmamış'}<br>
                    ${matrix.paar_program ? `<details><summary>Program</summary><pre>${JSON.stringify(JSON.parse(matrix.paar_program), null, 2)}</pre></details>` : ''}
                </div>
                
                <div class="algorithm-result result-slp mb-3">
                    <strong>SLP Heuristic:</strong><br>
                    XOR: ${matrix.slp_xor_count || 'Hesaplanmamış'}<br>
                    ${matrix.slp_program ? `<details><summary>Program</summary><pre>${JSON.stringify(JSON.parse(matrix.slp_program), null, 2)}</pre></details>` : ''}
                </div>
                
                <div class="mt-3">
                    <small class="text-muted">
                        <strong>Oluşturulma:</strong> ${formatDate(matrix.created_at)}<br>
                        <strong>Son Güncelleme:</strong> ${formatDate(matrix.updated_at)}
                    </small>
                </div>
            </div>
        </div>
    `;
    
    modalBody.innerHTML = html;
}

// Recalculate matrix algorithms
async function recalculateMatrix() {
    if (!currentMatrixId) return;
    
    try {
        showLoading('Algoritmalar yeniden hesaplanıyor...');
        
        const response = await fetch('/api/matrices/recalculate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                matrix_id: currentMatrixId,
                algorithms: ['boyar', 'paar', 'slp']
            })
        });
        
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.message || 'Yeniden hesaplama başarısız');
        }
        
        showAlert('Algoritmalar başarıyla yeniden hesaplandı!', 'success');
        
        // Update modal content
        displayMatrixDetails(data);
        
        // Reload matrices list
        loadMatrices();
        
    } catch (error) {
        console.error('Error recalculating matrix:', error);
        showAlert('Yeniden hesaplama sırasında hata oluştu: ' + error.message, 'danger');
    } finally {
        hideLoading();
    }
}

// Toggle input method between text and file
function toggleInputMethod(method) {
    const textSection = document.getElementById('textInputSection');
    const fileSection = document.getElementById('fileInputSection');
    
    if (method === 'file') {
        textSection.style.display = 'none';
        fileSection.style.display = 'block';
    } else {
        textSection.style.display = 'block';
        fileSection.style.display = 'none';
    }
}

// Bulk upload matrices
async function bulkUploadMatrices() {
    try {
        const inputMethod = document.querySelector('input[name="inputMethod"]:checked').value;
        const processImmediately = document.getElementById('bulkProcessImmediately').checked;
        const matrixRowCount = parseInt(document.getElementById('matrixRowCount').value);
        
        // Chunk boyutunu otomatik hesapla: 10 matris bloğu
        const chunkSize = matrixRowCount * 10;
        
        let content = '';
        
        if (inputMethod === 'file') {
            const fileInput = document.getElementById('bulkFile');
            if (!fileInput.files.length) {
                throw new Error('Lütfen bir dosya seçin');
            }
            
            const file = fileInput.files[0];
            content = await readFileContent(file);
        } else {
            content = document.getElementById('bulkData').value;
            if (!content.trim()) {
                throw new Error('Lütfen matris verilerini girin');
            }
        }
        
        // Process content in chunks
        await processContentInChunks(content, chunkSize, matrixRowCount, processImmediately);
        
    } catch (error) {
        console.error('Error in bulk upload:', error);
        showAlert('Toplu yükleme sırasında hata oluştu: ' + error.message, 'danger');
    } finally {
        showBulkUploadProgress(false);
    }
}

// Read file content
function readFileContent(file) {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.onload = function(e) {
            resolve(e.target.result);
        };
        reader.onerror = function() {
            reject(new Error('Dosya okunamadı'));
        };
        reader.readAsText(file, 'utf-8');
    });
}

// Process content in chunks
async function processContentInChunks(content, chunkSize, matrixRowCount, processImmediately) {
    const lines = content.split('\n').map(line => line.trim()).filter(line => line);
    const totalLines = lines.length;
    
    if (totalLines === 0) {
        throw new Error('Dosya boş veya geçersiz');
    }
    
    // Calculate optimal chunk size based on matrix row count
    const adjustedChunkSize = Math.max(chunkSize, matrixRowCount * 2); // En az 2 matris içerecek şekilde
    const totalChunks = Math.ceil(totalLines / adjustedChunkSize);
    
    showBulkUploadProgress(true);
    updateBulkUploadProgress(0, totalChunks, 'Chunk işleme başlıyor...');
    updateChunkProgress(0, totalLines, 'Dosya analiz ediliyor...');
    
    const allResults = [];
    let processedLines = 0;
    let remainingContent = '';
    
    for (let chunkIndex = 0; chunkIndex < totalChunks; chunkIndex++) {
        const startLine = chunkIndex * adjustedChunkSize;
        const endLine = Math.min(startLine + adjustedChunkSize, totalLines);
        
        // Get chunk lines
        const chunkLines = lines.slice(startLine, endLine);
        
        // Combine with remaining content from previous chunk
        const chunkContent = remainingContent + '\n' + chunkLines.join('\n');
        
        updateChunkProgress(endLine, totalLines, `Chunk ${chunkIndex + 1}/${totalChunks} işleniyor...`);
        
        // Parse matrices from chunk, keeping incomplete matrices for next chunk
        const { matrices, remaining } = parseChunkData(chunkContent, matrixRowCount);
        remainingContent = remaining;
        
        if (matrices.length > 0) {
            updateBulkUploadProgress(chunkIndex, totalChunks, `Chunk ${chunkIndex + 1}: ${matrices.length} matris yükleniyor...`);
            
            // Upload matrices from this chunk
            const chunkResults = await uploadMatricesChunk(matrices, processImmediately);
            allResults.push(...chunkResults);
        }
        
        processedLines = endLine;
        updateBulkUploadProgress(chunkIndex + 1, totalChunks, `Chunk ${chunkIndex + 1}/${totalChunks} tamamlandı`);
        
        // Small delay to prevent UI blocking
        await new Promise(resolve => setTimeout(resolve, 10));
    }
    
    // Process any remaining content
    if (remainingContent.trim()) {
        const { matrices } = parseChunkData(remainingContent, matrixRowCount);
        if (matrices.length > 0) {
            const finalResults = await uploadMatricesChunk(matrices, processImmediately);
            allResults.push(...finalResults);
        }
    }
    
    // Show final results
    showBulkUploadResults(allResults);
    
    // Clear form if all successful
    const successCount = allResults.filter(r => r.success).length;
    if (successCount === allResults.length && allResults.length > 0) {
        document.getElementById('bulkUploadForm').reset();
        showAlert(`${successCount} matris başarıyla yüklendi!`, 'success');
        
        // Switch to matrices tab and reload
        const matricesTab = new bootstrap.Tab(document.getElementById('matrices-tab'));
        matricesTab.show();
        loadMatrices();
    } else if (allResults.length > 0) {
        showAlert(`${successCount}/${allResults.length} matris yüklendi. Detaylar için sonuçları kontrol edin.`, 'warning');
    } else {
        showAlert('Hiç matris bulunamadı veya işlenemedi.', 'warning');
    }
}

// Parse chunk data with matrix boundary awareness
function parseChunkData(content, matrixRowCount) {
    const matrices = [];
    let remaining = '';
    
    // Split by "------------------------------" separator (gerçek dosya formatı)
    const sections = content.split('------------------------------');
    
    for (let i = 0; i < sections.length; i++) {
        const section = sections[i].trim();
        if (!section) continue;
        
        const lines = section.split('\n').map(line => line.trim()).filter(line => line);
        if (lines.length < 2) {
            // Incomplete section, add to remaining
            if (i === sections.length - 1) {
                remaining = section;
            }
            continue;
        }
        
        // Extract title (first line)
        const title = lines[0];
        
        // Find matrix data (lines with brackets)
        const matrixLines = [];
        for (let j = 1; j < lines.length; j++) {
            const line = lines[j];
            if (line.startsWith('[') && line.endsWith(']')) {
                matrixLines.push(line);
            }
        }
        
        // Check if we have complete matrix
        if (matrixLines.length === 0) {
            if (i === sections.length - 1) {
                remaining = section;
            }
            continue;
        }
        
        // If this is the last section and matrix seems incomplete, save for next chunk
        if (i === sections.length - 1 && matrixLines.length < 12 && sections.length > 1) {
            remaining = section;
            continue;
        }
        
        // Parse matrix
        const matrix = [];
        for (const line of matrixLines) {
            // Remove brackets and split by spaces or commas
            const cleanLine = line.replace(/[\[\]]/g, '').trim();
            const elements = cleanLine.split(/[\s,]+/).filter(el => el);
            if (elements.length > 0) {
                matrix.push(elements);
            }
        }
        
        if (matrix.length > 0) {
            matrices.push({
                title: title,
                matrix: matrix
            });
        }
    }
    
    return { matrices, remaining };
}

// Upload a chunk of matrices
async function uploadMatricesChunk(matrices, processImmediately) {
    const results = [];
    const endpoint = processImmediately ? '/api/matrices/process' : '/api/matrices';
    
    for (let i = 0; i < matrices.length; i++) {
        const matrix = matrices[i];
        
        try {
            const response = await fetch(endpoint, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    title: matrix.title,
                    matrix: matrix.matrix
                })
            });
            
            const data = await response.json();
            
            if (response.ok) {
                results.push({
                    success: true,
                    title: matrix.title,
                    id: data.id,
                    message: 'Başarıyla yüklendi'
                });
            } else {
                results.push({
                    success: false,
                    title: matrix.title,
                    message: data.message || 'Yükleme başarısız'
                });
            }
        } catch (error) {
            results.push({
                success: false,
                title: matrix.title,
                message: error.message
            });
        }
    }
    
    return results;
}

// Show/hide bulk upload progress
function showBulkUploadProgress(show) {
    const progressDiv = document.getElementById('bulkUploadProgress');
    const resultsDiv = document.getElementById('bulkUploadResults');
    
    if (show) {
        progressDiv.style.display = 'block';
        resultsDiv.style.display = 'none';
    } else {
        progressDiv.style.display = 'none';
    }
}

// Update bulk upload progress
function updateBulkUploadProgress(current, total, status) {
    const progressBar = document.querySelector('#bulkUploadProgress .progress-bar');
    const statusText = document.getElementById('bulkUploadStatus');
    
    const percentage = total > 0 ? (current / total) * 100 : 0;
    progressBar.style.width = percentage + '%';
    progressBar.textContent = Math.round(percentage) + '%';
    statusText.textContent = status;
}

// Update chunk progress
function updateChunkProgress(current, total, status) {
    const chunkProgressBar = document.getElementById('chunkProgress');
    const chunkStatusText = document.getElementById('chunkStatus');
    
    if (chunkProgressBar && chunkStatusText) {
        const percentage = total > 0 ? (current / total) * 100 : 0;
        chunkProgressBar.style.width = percentage + '%';
        chunkProgressBar.textContent = Math.round(percentage) + '%';
        chunkStatusText.textContent = status;
    }
}

// Show bulk upload results
function showBulkUploadResults(results) {
    const resultsDiv = document.getElementById('bulkUploadResults');
    const resultsList = document.getElementById('bulkUploadResultsList');
    
    const html = results.map(result => `
        <div class="alert ${result.success ? 'alert-success' : 'alert-danger'} py-2">
            <strong>${escapeHtml(result.title)}:</strong> ${escapeHtml(result.message)}
            ${result.id ? ` (ID: ${result.id})` : ''}
        </div>
    `).join('');
    
    resultsList.innerHTML = html;
    resultsDiv.style.display = 'block';
}

// Utility functions
function showLoading(text = 'Yükleniyor...') {
    document.getElementById('loadingText').textContent = text;
    const modalElement = document.getElementById('loadingModal');
    
    // Manuel olarak modal'ı göster
    modalElement.style.display = 'block';
    modalElement.classList.add('show');
    modalElement.setAttribute('aria-modal', 'true');
    modalElement.setAttribute('role', 'dialog');
    modalElement.removeAttribute('aria-hidden');
    
    // Backdrop ekle
    let backdrop = document.querySelector('.modal-backdrop');
    if (!backdrop) {
        backdrop = document.createElement('div');
        backdrop.className = 'modal-backdrop fade show';
        document.body.appendChild(backdrop);
    }
    
    // Body'ye modal-open class'ı ekle
    document.body.classList.add('modal-open');
    document.body.style.overflow = 'hidden';
}

function hideLoading() {
    const modalElement = document.getElementById('loadingModal');
    
    // Manuel olarak modal'ı gizle
    modalElement.style.display = 'none';
    modalElement.classList.remove('show');
    modalElement.setAttribute('aria-hidden', 'true');
    modalElement.removeAttribute('aria-modal');
    modalElement.removeAttribute('role');
    
    // Backdrop'ı kaldır
    const backdrop = document.querySelector('.modal-backdrop');
    if (backdrop) {
        backdrop.remove();
    }
    
    // Body'den modal-open class'ını kaldır
    document.body.classList.remove('modal-open');
    document.body.style.overflow = '';
    document.body.style.paddingRight = '';
    
    // Bootstrap modal instance varsa temizle
    const modal = bootstrap.Modal.getInstance(modalElement);
    if (modal) {
        modal.dispose();
    }
}

function showAlert(message, type = 'info') {
    // Create alert element
    const alertDiv = document.createElement('div');
    alertDiv.className = `alert alert-${type} alert-dismissible fade show position-fixed`;
    alertDiv.style.cssText = 'top: 20px; right: 20px; z-index: 9999; max-width: 400px;';
    alertDiv.innerHTML = `
        ${message}
        <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
    `;
    
    document.body.appendChild(alertDiv);
    
    // Auto remove after 5 seconds
    setTimeout(() => {
        if (alertDiv.parentNode) {
            alertDiv.remove();
        }
    }, 5000);
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleString('tr-TR');
}

// Debug functions - console'da kullanmak için
window.testShowLoading = function() {
    console.log('Testing showLoading...');
    showLoading('Test mesajı');
};

window.testHideLoading = function() {
    console.log('Testing hideLoading...');
    hideLoading();
};

window.checkModalState = function() {
    const modal = document.getElementById('loadingModal');
    const backdrop = document.querySelector('.modal-backdrop');
    console.log('Modal state:', {
        display: modal.style.display,
        classList: Array.from(modal.classList),
        ariaHidden: modal.getAttribute('aria-hidden'),
        ariaModal: modal.getAttribute('aria-modal'),
        role: modal.getAttribute('role'),
        backdropExists: !!backdrop,
        bodyHasModalOpen: document.body.classList.contains('modal-open'),
        bodyOverflow: document.body.style.overflow
    });
};

// Filter functions
function applyFilters() {
    // Get filter values
    const hamXorMin = document.getElementById('hamXorMin').value;
    const hamXorMax = document.getElementById('hamXorMax').value;
    const boyarXorMin = document.getElementById('boyarXorMin').value;
    const boyarXorMax = document.getElementById('boyarXorMax').value;
    const paarXorMin = document.getElementById('paarXorMin').value;
    const paarXorMax = document.getElementById('paarXorMax').value;
    const slpXorMin = document.getElementById('slpXorMin').value;
    const slpXorMax = document.getElementById('slpXorMax').value;

    // Update current filters
    currentFilters.hamXorMin = hamXorMin ? parseInt(hamXorMin) : null;
    currentFilters.hamXorMax = hamXorMax ? parseInt(hamXorMax) : null;
    currentFilters.boyarXorMin = boyarXorMin ? parseInt(boyarXorMin) : null;
    currentFilters.boyarXorMax = boyarXorMax ? parseInt(boyarXorMax) : null;
    currentFilters.paarXorMin = paarXorMin ? parseInt(paarXorMin) : null;
    currentFilters.paarXorMax = paarXorMax ? parseInt(paarXorMax) : null;
    currentFilters.slpXorMin = slpXorMin ? parseInt(slpXorMin) : null;
    currentFilters.slpXorMax = slpXorMax ? parseInt(slpXorMax) : null;

    // Reset to first page and reload
    currentPage = 1;
    loadMatrices();

    // Collapse the filter panel
    const filterCollapse = document.getElementById('filterCollapse');
    const bsCollapse = bootstrap.Collapse.getInstance(filterCollapse);
    if (bsCollapse) {
        bsCollapse.hide();
    }
}

function clearFilters() {
    // Clear all filter inputs
    document.getElementById('hamXorMin').value = '';
    document.getElementById('hamXorMax').value = '';
    document.getElementById('boyarXorMin').value = '';
    document.getElementById('boyarXorMax').value = '';
    document.getElementById('paarXorMin').value = '';
    document.getElementById('paarXorMax').value = '';
    document.getElementById('slpXorMin').value = '';
    document.getElementById('slpXorMax').value = '';

    // Clear current filters
    currentFilters = {
        hamXorMin: null,
        hamXorMax: null,
        boyarXorMin: null,
        boyarXorMax: null,
        paarXorMin: null,
        paarXorMax: null,
        slpXorMin: null,
        slpXorMax: null
    };

    // Reset to first page and reload
    currentPage = 1;
    loadMatrices();
}

// Bulk recalculate function
async function bulkRecalculate() {
    if (!confirm('Algoritma hesaplanmamış tüm matrisler için algoritmaları çalıştırmak istediğinizden emin misiniz? Bu işlem uzun sürebilir.')) {
        return;
    }

    try {
        showLoading('Toplu algoritma hesaplama başlatılıyor...');
        
        const response = await fetch('/api/matrices/bulk-recalculate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                algorithms: ['boyar', 'paar', 'slp'],
                limit: 100
            })
        });
        
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.message || 'Toplu hesaplama başlatılamadı');
        }
        
        showAlert(data.message, 'success');
        
        // Reload matrices list after a short delay
        setTimeout(() => {
            loadMatrices();
        }, 2000);
        
    } catch (error) {
        console.error('Error in bulk recalculate:', error);
        showAlert('Toplu hesaplama sırasında hata oluştu: ' + error.message, 'danger');
    } finally {
        hideLoading();
    }
}

// Calculate matrix inverse
async function calculateMatrixInverse() {
    console.log('calculateMatrixInverse called, currentMatrixId:', currentMatrixId);
    
    if (!currentMatrixId) {
        console.error('currentMatrixId is null or undefined');
        showAlert('Matris ID bulunamadı. Lütfen önce bir matris seçin.', 'warning');
        return;
    }
    
    try {
        showLoading('Ters matris hesaplanıyor...');
        
        console.log('Making API request to:', `/api/matrices/${currentMatrixId}/inverse`);
        
        const response = await fetch(`/api/matrices/${currentMatrixId}/inverse`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        });
        
        console.log('API response status:', response.status);
        
        const data = await response.json();
        
        console.log('API response data:', data);
        
        if (!response.ok) {
            throw new Error(data.message || 'Ters matris hesaplama başarısız');
        }
        
        showAlert(`Ters matris başarıyla hesaplandı ve kaydedildi! Yeni matris ID: ${data.id}`, 'success');
        
        // Reload matrices list to show the new inverse matrix
        loadMatrices();
        
        // Optionally, show the inverse matrix details
        setTimeout(() => {
            viewMatrix(data.id);
        }, 1000);
        
    } catch (error) {
        console.error('Error calculating matrix inverse:', error);
        showAlert('Ters matris hesaplama sırasında hata oluştu: ' + error.message, 'danger');
    } finally {
        hideLoading();
    }
}

// Bulk Inverse Functions
async function previewBulkInverse() {
    const maxSmallestXor = document.getElementById('maxSmallestXor').value;
    const skipExisting = document.getElementById('skipExisting').checked;
    
    if (!maxSmallestXor) {
        showAlert('Lütfen maksimum smallest XOR değerini girin', 'warning');
        return;
    }

    showLoading('Önizleme hazırlanıyor...');
    
    try {
        // Get all matrices without any XOR filter - we'll filter on smallest_xor field
        const response = await fetch(`/api/matrices?limit=10000`);
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.message || 'Matrisler alınamadı');
        }

        const matrices = data.matrices || [];
        const candidateMatrices = matrices.filter(matrix => {
            // Filter based on smallest_xor field specifically
            const smallestXor = matrix.smallest_xor;
            if (!smallestXor || smallestXor >= parseInt(maxSmallestXor)) {
                return false;
            }
            
            if (skipExisting && matrix.inverse_matrix_id) {
                return false;
            }
            
            return true;
        });

        // Update stats
        document.getElementById('totalMatricesCount').textContent = matrices.length;
        document.getElementById('withInverseCount').textContent = 
            matrices.filter(m => m.inverse_matrix_id).length;

        // Show preview
        const previewDiv = document.getElementById('bulkInversePreview');
        const previewContent = document.getElementById('previewContent');
        
        if (candidateMatrices.length === 0) {
            previewContent.innerHTML = `
                <div class="alert alert-warning">
                    <i class="fas fa-exclamation-triangle me-2"></i>
                    Belirtilen kriterlere uygun matris bulunamadı.
                </div>
            `;
        } else {
            let html = `
                <div class="alert alert-success">
                    <i class="fas fa-check-circle me-2"></i>
                    <strong>${candidateMatrices.length}</strong> matris için ters matris hesaplanacak.
                </div>
                <div class="table-responsive">
                    <table class="table table-sm">
                        <thead>
                            <tr>
                                <th>ID</th>
                                <th>Başlık</th>
                                <th>Grup</th>
                                <th>Smallest XOR</th>
                                <th>Ters Matris</th>
                            </tr>
                        </thead>
                        <tbody>
            `;
            
            candidateMatrices.slice(0, 20).forEach(matrix => {
                html += `
                    <tr>
                        <td>${matrix.id}</td>
                        <td>${escapeHtml(matrix.title)}</td>
                        <td>${escapeHtml(matrix.group || '-')}</td>
                        <td>${matrix.smallest_xor || '-'}</td>
                        <td>
                            ${matrix.inverse_matrix_id ? 
                                '<span class="badge bg-success">Var</span>' : 
                                '<span class="badge bg-warning">Yok</span>'
                            }
                        </td>
                    </tr>
                `;
            });
            
            html += '</tbody></table>';
            
            if (candidateMatrices.length > 20) {
                html += `
                    <div class="text-muted text-center mt-2">
                        <small>... ve ${candidateMatrices.length - 20} matris daha</small>
                    </div>
                `;
            }
            
            html += '</div>';
            previewContent.innerHTML = html;
        }
        
        previewDiv.style.display = 'block';
        
    } catch (error) {
        console.error('Preview error:', error);
        showAlert('Önizleme hazırlanırken hata oluştu: ' + error.message, 'danger');
    } finally {
        hideLoading();
    }
}

async function startBulkInverse() {
    const maxSmallestXor = document.getElementById('maxSmallestXor').value;
    const batchSize = parseInt(document.getElementById('batchSize').value) || 10;
    const skipExisting = document.getElementById('skipExisting').checked;
    const calculateAlgorithms = document.getElementById('calculateAlgorithms').checked;
    
    if (!maxSmallestXor) {
        showAlert('Lütfen maksimum smallest XOR değerini girin', 'warning');
        return;
    }

    // Hide preview and results, show progress
    document.getElementById('bulkInversePreview').style.display = 'none';
    document.getElementById('bulkInverseResults').style.display = 'none';
    document.getElementById('bulkInverseProgress').style.display = 'block';
    
    updateBulkInverseProgress(0, 0, 'İşlem başlatılıyor...');
    
    try {
        const response = await fetch('/api/matrices/bulk-inverse', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                max_smallest_xor: parseInt(maxSmallestXor),
                batch_size: batchSize,
                skip_existing: skipExisting,
                calculate_algorithms: calculateAlgorithms
            })
        });
        
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.message || 'Toplu ters alma başlatılamadı');
        }
        
        updateBulkInverseProgress(0, 0, data.message);
        
        // Show success message
        setTimeout(() => {
            showBulkInverseResults([]);
            showAlert('Toplu ters alma işlemi arka planda başlatıldı. Sonuçları görmek için matrisler sayfasını yenileyin.', 'success');
        }, 2000);
        
    } catch (error) {
        console.error('Bulk inverse error:', error);
        updateBulkInverseProgress(0, 0, 'Hata oluştu: ' + error.message);
        showAlert('Toplu ters alma işleminde hata oluştu: ' + error.message, 'danger');
    }
}

function updateBulkInverseProgress(current, total, status) {
    const progressBar = document.getElementById('inverseProgressBar');
    const progressText = document.getElementById('inverseProgressText');
    const statusText = document.getElementById('inverseStatusText');
    
    const percentage = total > 0 ? Math.round((current / total) * 100) : 0;
    
    progressBar.style.width = percentage + '%';
    progressBar.textContent = percentage + '%';
    progressText.textContent = `${current} / ${total}`;
    statusText.textContent = status;
}

function showBulkInverseResults(results) {
    const resultsDiv = document.getElementById('bulkInverseResults');
    const resultsContent = document.getElementById('resultsContent');
    
    let html = `
        <div class="alert alert-info">
            <i class="fas fa-info-circle me-2"></i>
            Toplu ters alma işlemi arka planda çalışmaktadır. Sonuçları görmek için matrisler sayfasını kontrol edin.
        </div>
    `;
    
    resultsContent.innerHTML = html;
    resultsDiv.style.display = 'block';
    
    // Hide progress
    document.getElementById('bulkInverseProgress').style.display = 'none';
}

async function loadBulkInverseStats() {
    try {
        const response = await fetch('/api/matrices?limit=1');
        const data = await response.json();
        
        if (response.ok) {
            document.getElementById('totalMatricesCount').textContent = data.total || 0;
            
            // Count matrices with inverse
            const allResponse = await fetch('/api/matrices?limit=1000');
            const allData = await allResponse.json();
            
            if (allResponse.ok) {
                const withInverse = (allData.matrices || []).filter(m => m.inverse_matrix_id).length;
                document.getElementById('withInverseCount').textContent = withInverse;
            }
        }
    } catch (error) {
        console.error('Stats loading error:', error);
    }
}

// Load inverse pairs
async function loadInversePairs() {
    try {
        showPairsLoading(true);
        
        const params = new URLSearchParams({
            page: currentPairPage,
            limit: currentPairLimit
        });
        
        if (currentPairGroupFilter) {
            params.append('group', currentPairGroupFilter);
        }
        
        if (currentMaxCombinedXor) {
            params.append('max_combined_xor', currentMaxCombinedXor);
        }
        
        params.append('sort', currentSortOrder);
        
        const response = await fetch(`/api/inverse-pairs?${params}`);
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.message || 'Ters matris çiftleri alınamadı');
        }
        
        displayInversePairs(data.pairs || []);
        setupPairsPagination(data.page, data.total_pages, data.total);
        
    } catch (error) {
        console.error('Error loading inverse pairs:', error);
        showAlert('Ters matris çiftleri yüklenirken hata oluştu: ' + error.message, 'danger');
        document.getElementById('inversePairsContainer').innerHTML = `
            <div class="alert alert-danger">
                <i class="fas fa-exclamation-triangle me-2"></i>
                Ters matris çiftleri yüklenemedi: ${error.message}
            </div>
        `;
    } finally {
        showPairsLoading(false);
    }
}

// Display inverse pairs
function displayInversePairs(pairs) {
    const container = document.getElementById('inversePairsContainer');
    
    if (!pairs || pairs.length === 0) {
        container.innerHTML = `
            <div class="text-center py-5">
                <i class="fas fa-link fa-3x text-muted mb-3"></i>
                <h5 class="text-muted">Henüz ters matris çifti bulunmuyor</h5>
                <p class="text-muted">Ters matris hesaplamak için "Toplu Ters Alma" sekmesini kullanın.</p>
            </div>
        `;
        return;
    }

    const html = pairs.map(pair => `
        <div class="card mb-3">
            <div class="card-header">
                <div class="row align-items-center">
                    <div class="col-md-8">
                        <h6 class="mb-0">
                            <i class="fas fa-link me-2 text-primary"></i>
                            Matris Çifti - Toplam XOR: <span class="badge bg-primary">${pair.combined_xor}</span>
                        </h6>
                    </div>
                    <div class="col-md-4 text-end">
                        <small class="text-muted">Grup: ${pair.group || 'Belirtilmemiş'}</small>
                    </div>
                </div>
            </div>
            <div class="card-body">
                <div class="row">
                    <!-- Original Matrix -->
                    <div class="col-md-6">
                        <div class="border rounded p-3 h-100">
                            <div class="d-flex justify-content-between align-items-center mb-2">
                                <h6 class="text-success mb-0">
                                    <i class="fas fa-table me-1"></i>Orijinal Matris
                                </h6>
                                <button class="btn btn-sm btn-outline-success" onclick="viewMatrix(${pair.original_id})">
                                    <i class="fas fa-eye me-1"></i>Detay
                                </button>
                            </div>
                            <div class="mb-2">
                                <strong>ID:</strong> ${pair.original_id} | 
                                <strong>Başlık:</strong> ${escapeHtml(pair.original_title)}
                            </div>
                            <div class="mb-2">
                                <strong>Smallest XOR:</strong> <span class="badge bg-success">${pair.original_xor}</span>
                            </div>
                            <div class="row text-center">
                                <div class="col-4">
                                    <small class="text-muted">Boyar</small><br>
                                    <span class="badge bg-secondary">${pair.original_boyar_xor || 'N/A'}</span>
                                </div>
                                <div class="col-4">
                                    <small class="text-muted">Paar</small><br>
                                    <span class="badge bg-secondary">${pair.original_paar_xor || 'N/A'}</span>
                                </div>
                                <div class="col-4">
                                    <small class="text-muted">SLP</small><br>
                                    <span class="badge bg-secondary">${pair.original_slp_xor || 'N/A'}</span>
                                </div>
                            </div>
                        </div>
                    </div>
                    
                    <!-- Inverse Matrix -->
                    <div class="col-md-6">
                        <div class="border rounded p-3 h-100">
                            <div class="d-flex justify-content-between align-items-center mb-2">
                                <h6 class="text-warning mb-0">
                                    <i class="fas fa-exchange-alt me-1"></i>Ters Matris
                                </h6>
                                <button class="btn btn-sm btn-outline-warning" onclick="viewMatrix(${pair.inverse_id})">
                                    <i class="fas fa-eye me-1"></i>Detay
                                </button>
                            </div>
                            <div class="mb-2">
                                <strong>ID:</strong> ${pair.inverse_id} | 
                                <strong>Başlık:</strong> ${escapeHtml(pair.inverse_title)}
                            </div>
                            <div class="mb-2">
                                <strong>Smallest XOR:</strong> <span class="badge bg-warning">${pair.inverse_xor}</span>
                            </div>
                            <div class="row text-center">
                                <div class="col-4">
                                    <small class="text-muted">Boyar</small><br>
                                    <span class="badge bg-secondary">${pair.inverse_boyar_xor || 'N/A'}</span>
                                </div>
                                <div class="col-4">
                                    <small class="text-muted">Paar</small><br>
                                    <span class="badge bg-secondary">${pair.inverse_paar_xor || 'N/A'}</span>
                                </div>
                                <div class="col-4">
                                    <small class="text-muted">SLP</small><br>
                                    <span class="badge bg-secondary">${pair.inverse_slp_xor || 'N/A'}</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                
                <!-- Additional Info -->
                <div class="row mt-3">
                    <div class="col-12">
                        <div class="bg-light rounded p-2">
                            <div class="row text-center">
                                <div class="col-md-3">
                                    <small class="text-muted">Toplam XOR</small><br>
                                    <strong class="text-primary">${pair.combined_xor}</strong>
                                </div>
                                <div class="col-md-3">
                                    <small class="text-muted">XOR Farkı</small><br>
                                    <strong>${Math.abs(pair.original_xor - pair.inverse_xor)}</strong>
                                </div>
                                <div class="col-md-3">
                                    <small class="text-muted">Oluşturulma</small><br>
                                    <small>${formatDate(pair.created_at)}</small>
                                </div>
                                <div class="col-md-3">
                                    <small class="text-muted">Güncelleme</small><br>
                                    <small>${formatDate(pair.updated_at)}</small>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    `).join('');

    container.innerHTML = html;
}

// Setup pairs pagination
function setupPairsPagination(currentPage, totalPages, totalItems) {
    const container = document.getElementById('pairsPaginationContainer');
    const pagination = document.getElementById('pairsPagination');
    const paginationInfo = document.getElementById('pairsPaginationInfo');
    
    if (totalPages <= 1) {
        container.style.display = 'none';
        return;
    }
    
    container.style.display = 'block';
    
    // Update pagination info
    const startItem = (currentPage - 1) * currentPairLimit + 1;
    const endItem = Math.min(currentPage * currentPairLimit, totalItems);
    paginationInfo.textContent = `${startItem}-${endItem} / ${totalItems} çift gösteriliyor (Sayfa ${currentPage}/${totalPages})`;
    
    let html = '';
    
    // Previous button
    html += `
        <li class="page-item ${currentPage === 1 ? 'disabled' : ''}">
            <a class="page-link" href="#" ${currentPage > 1 ? `onclick="changePairPage(${currentPage - 1})"` : ''} aria-label="Previous">
                <span aria-hidden="true">&laquo;</span>
            </a>
        </li>
    `;
    
    // Page numbers
    const startPage = Math.max(1, currentPage - 2);
    const endPage = Math.min(totalPages, currentPage + 2);
    
    if (startPage > 1) {
        html += `<li class="page-item"><a class="page-link" href="#" onclick="changePairPage(1)">1</a></li>`;
        if (startPage > 2) {
            html += `<li class="page-item disabled"><span class="page-link">...</span></li>`;
        }
    }
    
    for (let i = startPage; i <= endPage; i++) {
        html += `
            <li class="page-item ${i === currentPage ? 'active' : ''}">
                <a class="page-link" href="#" ${i !== currentPage ? `onclick="changePairPage(${i})"` : ''}>${i}</a>
            </li>
        `;
    }
    
    if (endPage < totalPages) {
        if (endPage < totalPages - 1) {
            html += `<li class="page-item disabled"><span class="page-link">...</span></li>`;
        }
        html += `<li class="page-item"><a class="page-link" href="#" onclick="changePairPage(${totalPages})">${totalPages}</a></li>`;
    }
    
    // Next button
    html += `
        <li class="page-item ${currentPage === totalPages ? 'disabled' : ''}">
            <a class="page-link" href="#" ${currentPage < totalPages ? `onclick="changePairPage(${currentPage + 1})"` : ''} aria-label="Next">
                <span aria-hidden="true">&raquo;</span>
            </a>
        </li>
    `;
    
    pagination.innerHTML = html;
}

// Change pair page
function changePairPage(page) {
    if (page < 1) return;
    currentPairPage = page;
    loadInversePairs();
}

// Apply pair filters
function applyPairFilters() {
    currentPairGroupFilter = document.getElementById('pairGroupFilter').value.trim();
    const maxCombined = document.getElementById('maxCombinedXor').value;
    currentMaxCombinedXor = maxCombined ? parseInt(maxCombined) : null;
    currentPairLimit = parseInt(document.getElementById('pairPageSize').value);
    currentPairPage = 1; // Reset to first page
    loadInversePairs();
}

// Set sort order
function setSortOrder(order) {
    currentSortOrder = order;
    currentPairPage = 1; // Reset to first page
    loadInversePairs();
}

// Show/hide pairs loading
function showPairsLoading(show) {
    const loadingDiv = document.getElementById('pairsLoading');
    const containerDiv = document.getElementById('inversePairsContainer');
    const paginationDiv = document.getElementById('pairsPaginationContainer');
    
    if (show) {
        loadingDiv.style.display = 'block';
        containerDiv.style.display = 'none';
        paginationDiv.style.display = 'none';
    } else {
        loadingDiv.style.display = 'none';
        containerDiv.style.display = 'block';
    }
} 
