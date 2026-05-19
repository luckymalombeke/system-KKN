import React from 'react';
import { motion } from 'framer-motion';
import { X, Mail, CreditCard, MapPin } from 'lucide-react';

const PanduanModal = ({ isOpen, onClose }) => {
  if (!isOpen) return null;

  return (
    <div style={{
      position: 'fixed', top: 0, left: 0, right: 0, bottom: 0,
      background: 'rgba(2, 6, 23, 0.8)',
      backdropFilter: 'blur(8px)',
      display: 'flex', justifyContent: 'center', alignItems: 'center',
      zIndex: 100
    }}>
      <motion.div 
        initial={{ scale: 0.9, opacity: 0, y: 20 }}
        animate={{ scale: 1, opacity: 1, y: 0 }}
        exit={{ scale: 0.9, opacity: 0, y: 20 }}
        className="glass-card"
        style={{ width: '100%', maxWidth: '500px', position: 'relative', padding: '2.5rem' }}
      >
        <button 
          onClick={onClose}
          style={{ position: 'absolute', top: '15px', right: '15px', background: 'transparent', border: 'none', color: 'var(--text-dim)', cursor: 'pointer' }}
        >
          <X size={24} />
        </button>

        <div style={{ textAlign: 'center', marginBottom: '2rem' }}>
          <h2 className="gradient-text" style={{ fontSize: '1.8rem', marginBottom: '0.5rem' }}>
            Panduan Sistem KKN
          </h2>
          <p style={{ color: 'var(--text-dim)', fontSize: '0.95rem' }}>
            Ikuti 3 langkah mudah berikut untuk menyelesaikan administrasi KKN Anda.
          </p>
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
          {/* Langkah 1 */}
          <div style={{ display: 'flex', gap: '1rem', alignItems: 'flex-start' }}>
            <div style={{ padding: '10px', background: 'rgba(34, 211, 238, 0.1)', borderRadius: '12px', marginTop: '4px' }}>
              <Mail color="var(--primary)" size={20} />
            </div>
            <div>
              <h4 style={{ fontSize: '1.1rem', marginBottom: '0.2rem' }}>1. Login Tanpa Password</h4>
              <p style={{ color: 'var(--text-dim)', fontSize: '0.9rem', lineHeight: '1.5' }}>
                Klik "Mulai Pendaftaran", ketik NIM Anda. Kami akan mengirimkan 6-digit kode keamanan (OTP) ke email kampus Anda. Tidak perlu pusing mengingat password!
              </p>
            </div>
          </div>

          {/* Langkah 2 */}
          <div style={{ display: 'flex', gap: '1rem', alignItems: 'flex-start' }}>
            <div style={{ padding: '10px', background: 'rgba(234, 179, 8, 0.1)', borderRadius: '12px', marginTop: '4px' }}>
              <CreditCard color="#eab308" size={20} />
            </div>
            <div>
              <h4 style={{ fontSize: '1.1rem', marginBottom: '0.2rem' }}>2. Bayar Tagihan (Midtrans)</h4>
              <p style={{ color: 'var(--text-dim)', fontSize: '0.9rem', lineHeight: '1.5' }}>
                Setelah masuk ke Dashboard, klik "Bayar Sekarang". Anda bisa membayar menggunakan GoPay, Virtual Account BCA, atau minimarket. Status lunas akan otomatis tercatat.
              </p>
            </div>
          </div>

          {/* Langkah 3 */}
          <div style={{ display: 'flex', gap: '1rem', alignItems: 'flex-start' }}>
            <div style={{ padding: '10px', background: 'rgba(239, 68, 68, 0.1)', borderRadius: '12px', marginTop: '4px' }}>
              <MapPin color="#ef4444" size={20} />
            </div>
            <div>
              <h4 style={{ fontSize: '1.1rem', marginBottom: '0.2rem' }}>3. Pantau Plotting Lokasi</h4>
              <p style={{ color: 'var(--text-dim)', fontSize: '0.9rem', lineHeight: '1.5' }}>
                Cek Dashboard secara berkala. Setelah Anda lunas, Admin akan menempatkan Anda ke Desa KKN. Lokasi dan info penting akan muncul di Papan Notifikasi Anda.
              </p>
            </div>
          </div>
        </div>

        <button onClick={onClose} className="btn-primary" style={{ width: '100%', justifyContent: 'center', marginTop: '2.5rem' }}>
          Saya Mengerti
        </button>

      </motion.div>
    </div>
  );
};

export default PanduanModal;
