const BASE_URL = "http://localhost:8081/api/v1";

const getToken = () => localStorage.getItem("kkn_token");

const authHeaders = (extra = {}) => {
    const headers = { "Content-Type": "application/json", ...extra };
    const token = getToken();
    if (token) {
        headers.Authorization = `Bearer ${token}`;
    }
    return headers;
};

const authFetch = (url, options = {}) => {
    const token = getToken();
    if (!token) {
        return Promise.reject(new Error("Sesi login tidak ditemukan, silakan login kembali"));
    }
    return fetch(url, {
        ...options,
        headers: {
            ...authHeaders(),
            ...(options.headers || {}),
        },
    });
};

const handleResponse = async (response) => {
    if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || errorData.message || "Terjadi kesalahan jaringan/server");
    }
    return response.json();
};

export const KKN_API = {

    // ==========================================
    // AUTHENTICATION (PASSWORDLESS OTP)
    // ==========================================
    requestOtp: async (nim) => {
        const response = await fetch(`${BASE_URL}/auth/request-otp`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ nim }),
        });
        return handleResponse(response);
    },

    verifyOtp: async (nim, otp) => {
        const response = await fetch(`${BASE_URL}/auth/verify-otp`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ nim, otp }),
        });
        return handleResponse(response);
    },

    adminLogin: async (email, password) => {
        const response = await fetch(`${BASE_URL}/auth/admin-login`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ email, password }),
        });
        return handleResponse(response);
    },

    // ==========================================
    // PESERTA (Publik)
    // ==========================================
    registerPeserta: async (data) => {
        const response = await fetch(`${BASE_URL}/peserta/register`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data),
        });
        return handleResponse(response);
    },

    // ==========================================
    // MAHASISWA (JWT — user_id dari token)
    // ==========================================
    getMyProfile: async () => {
        const response = await authFetch(`${BASE_URL}/mahasiswa/profile`);
        return handleResponse(response);
    },

    getMyPembayaran: async () => {
        const response = await authFetch(`${BASE_URL}/mahasiswa/pembayaran`);
        return handleResponse(response);
    },

    createMyInvoice: async (amount) => {
        const response = await authFetch(`${BASE_URL}/mahasiswa/pembayaran/invoice`, {
            method: "POST",
            body: JSON.stringify({ amount: parseInt(amount, 10) }),
        });
        return handleResponse(response);
    },

    getMyNotifikasi: async () => {
        const response = await authFetch(`${BASE_URL}/mahasiswa/notifikasi`);
        return handleResponse(response);
    },

    markMyNotifikasiRead: async (notifikasiId) => {
        const response = await authFetch(`${BASE_URL}/mahasiswa/notifikasi/${notifikasiId}/read`, {
            method: "PATCH",
        });
        return handleResponse(response);
    },

    // ==========================================
    // WEBHOOK (Publik — simulasi Midtrans)
    // ==========================================
    simulateWebhook: async (orderId, status, grossAmount = "500000.00") => {
        const response = await fetch(`${BASE_URL}/pembayaran/notification`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                order_id: orderId,
                status_code: "200",
                gross_amount: grossAmount,
                transaction_status: status,
                fraud_status: "accept",
            }),
        });
        return handleResponse(response);
    },

    // ==========================================
    // ADMIN (JWT admin)
    // ==========================================
    getAllPeserta: async () => {
        const response = await authFetch(`${BASE_URL}/admin/peserta`);
        return handleResponse(response);
    },

    getPesertaById: async (id) => {
        const response = await authFetch(`${BASE_URL}/admin/peserta/${id}`);
        return handleResponse(response);
    },

    updatePesertaStatus: async (id, status) => {
        const response = await authFetch(`${BASE_URL}/admin/peserta/${id}/status`, {
            method: "PATCH",
            body: JSON.stringify({ status }),
        });
        return handleResponse(response);
    },

    assignLokasi: async (pesertaId, lokasiId) => {
        const response = await authFetch(`${BASE_URL}/admin/peserta/${pesertaId}/lokasi`, {
            method: "PATCH",
            body: JSON.stringify({ lokasi_id: lokasiId }),
        });
        return handleResponse(response);
    },

    getLokasi: async () => {
        const response = await authFetch(`${BASE_URL}/admin/lokasi`);
        return handleResponse(response);
    },

    createLokasi: async (data) => {
        const response = await authFetch(`${BASE_URL}/admin/lokasi`, {
            method: "POST",
            body: JSON.stringify(data),
        });
        return handleResponse(response);
    },
};
