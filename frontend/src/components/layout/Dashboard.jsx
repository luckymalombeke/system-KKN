import React from 'react';
import { motion } from 'framer-motion';
import { 
  User, 
  MapPin, 
  CreditCard, 
  Bell, 
  LogOut,
  ChevronRight,
  ShieldCheck,
  AlertCircle,
  RefreshCw
} from 'lucide-react';

import { KKN_API } from '../../services/api';

const formatWIB = (isoString) => {
  if (!isoString) return '-';
  return new Date(isoString).toLocaleString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    timeZone: 'Asia/Jakarta',
  }) + ' WIB';
};

const formatRupiah = (amount) =>
  new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(amount || 0);

const Dashboard = ({ userNim, onLogout }) => {
  const [payment, setPayment] = React.useState(null);
  const [isLoading, setIsLoading] = React.useState(false);
  const [profile, setProfile] = React.useState(null);
  const [notifikasi, setNotifikasi] = React.useState([]);

  const paymentStatus = payment?.status ?? 'none';
  const canCreateInvoice = paymentStatus === 'none' || paymentStatus === 'failed';

  const loadPayment = React.useCallback(async () => {
    try {
      const payRes = await KKN_API.getMyPembayaran();
      setPayment(payRes.data);
    } catch {
      setPayment(null);
    }
  }, []);

  React.useEffect(() => {
    const loadDashboardData = async () => {
      try {
        const profileRes = await KKN_API.getMyProfile();
        setProfile(profileRes.data);
      } catch {
        setProfile(null);
      }

      await loadPayment();

      try {
        const notifRes = await KKN_API.getMyNotifikasi();
        setNotifikasi(notifRes.data || []);
      } catch {
        setNotifikasi([]);
      }
    };

    if (localStorage.getItem('kkn_token')) {
      loadDashboardData();
    }
  }, [loadPayment]);

  const handleGenerateInvoice = async () => {
    setIsLoading(true);
    try {
      const res = await KKN_API.createMyInvoice(500000);
      setPayment(res.data);
      if (res.data.payment_url) {
        window.open(res.data.payment_url, '_blank', 'noopener,noreferrer');
      }
    } catch (err) {
      alert('Gagal membuat tagihan: ' + err.message);
    } finally {
      setIsLoading(false);
    }
  };

  const handleOpenPayment = () => {
    if (payment?.payment_url) {
      window.open(payment.payment_url, '_blank', 'noopener,noreferrer');
    }
  };

  const handleSimulasiMidtrans = async () => {
    if (!payment?.external_id) return;
    setIsLoading(true);
    try {
      const amount = `${payment.amount}.00`;
      await KKN_API.simulateWebhook(payment.external_id, 'settlement', amount);
      await loadPayment();
    } catch (err) {
      alert('Gagal update Midtrans: ' + err.message);
    } finally {
      setIsLoading(false);
    }
  };

  const handleSimulasiExpire = async () => {
    if (!payment?.external_id) return;
    setIsLoading(true);
    try {
      const amount = `${payment.amount}.00`;
      await KKN_API.simulateWebhook(payment.external_id, 'expire', amount);
      await loadPayment();
    } catch (err) {
      alert('Gagal simulasi kedaluwarsa: ' + err.message);
    } finally {
      setIsLoading(false);
    }
  };

  const paymentSubtitle = {
    none: 'Tagihan belum dibuat',
    pending: 'Menunggu pembayaran',
    success: 'Pembayaran berhasil',
    failed: 'Tagihan kedaluwarsa / dibatalkan',
  }[paymentStatus] || 'Status pembayaran';

  const cardAccent =
    paymentStatus === 'success' ? '#22c55e' : paymentStatus === 'failed' ? '#ef4444' : '#eab308';

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: { opacity: 1, transition: { staggerChildren: 0.15 } }
  };

  const itemVariants = {
    hidden: { y: 20, opacity: 0 },
    visible: { y: 0, opacity: 1, transition: { type: 'spring', stiffness: 100 } }
  };

  const statusLabel = {
    pending: { text: 'Menunggu Verifikasi', color: '#eab308' },
    approved: { text: 'Telah Disetujui', color: 'var(--primary)' },
    rejected: { text: 'Ditolak', color: '#ef4444' },
  };
  const biodataStatus = statusLabel[profile?.status] || statusLabel.pending;

  return (
    <motion.div style={{ minHeight: '100vh', padding: '2rem 5%', position: 'relative', zIndex: 1 }}>
      
      <motion.div 
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '3rem', flexWrap: 'wrap', gap: '1rem' }}
      >
        <div>
          <h1 style={{ fontSize: '2rem', fontWeight: 700 }}>
            Halo, <span className="gradient-text">{userNim}</span> 👋
          </h1>
          <p style={{ color: 'var(--text-dim)' }}>Berikut adalah status pendaftaran dan informasi KKN Anda.</p>
        </div>
        
        <button 
          onClick={onLogout} 
          className="btn-outline" 
          style={{ borderColor: 'rgba(239, 68, 68, 0.5)', color: '#ef4444', padding: '0.6rem 1rem' }}
        >
          <LogOut size={18} /> Keluar
        </button>
      </motion.div>

      <motion.div 
        variants={containerVariants}
        initial="hidden"
        animate="visible"
        style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))',
          gap: '2rem'
        }}
      >
        <motion.div variants={itemVariants} className="glass-card" style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
          <motion.div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
            <div style={{ padding: '12px', background: 'rgba(34, 211, 238, 0.1)', borderRadius: '12px' }}>
              <User color="var(--primary)" size={24} />
            </div>
            <div>
              <h3 style={{ fontSize: '1.1rem' }}>Status Biodata</h3>
              <p style={{ color: 'var(--text-dim)', fontSize: '0.9rem' }}>Verifikasi Sistem</p>
            </div>
          </motion.div>
          <div style={{ padding: '1rem', background: 'rgba(34, 211, 238, 0.05)', border: '1px solid var(--border)', borderRadius: '12px', display: 'flex', alignItems: 'center', gap: '0.8rem' }}>
            <ShieldCheck color={biodataStatus.color} size={20} />
            <span style={{ fontWeight: 600, color: biodataStatus.color }}>{biodataStatus.text}</span>
          </div>
          {profile?.nama && (
            <p style={{ color: 'var(--text-dim)', fontSize: '0.9rem' }}>{profile.nama} · {profile.prodi}</p>
          )}
        </motion.div>

        <motion.div variants={itemVariants} className="glass-card" style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
            <motion.div style={{ padding: '12px', background: `${cardAccent}22`, borderRadius: '12px' }}>
              <CreditCard color={cardAccent} size={24} />
            </motion.div>
            <div>
              <h3 style={{ fontSize: '1.1rem' }}>Biaya Pendaftaran</h3>
              <p style={{ color: 'var(--text-dim)', fontSize: '0.9rem' }}>{paymentSubtitle}</p>
            </div>
          </div>
          
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-end', marginTop: '0.5rem' }}>
            <div>
              <p style={{ fontSize: '0.85rem', color: 'var(--text-dim)' }}>Total Tagihan</p>
              <h2 style={{ fontSize: '1.8rem', fontWeight: 700, color: paymentStatus === 'success' ? '#22c55e' : 'inherit' }}>
                {formatRupiah(payment?.amount || 500000)}
              </h2>
            </div>
            {paymentStatus === 'pending' && payment?.expired_at && (
              <div style={{ textAlign: 'right' }}>
                <p style={{ fontSize: '0.8rem', color: 'var(--text-dim)' }}>Batas Bayar</p>
                <p style={{ fontSize: '0.85rem', fontWeight: 600, color: '#eab308' }}>{formatWIB(payment.expired_at)}</p>
              </div>
            )}
            {paymentStatus === 'success' && (
              <motion.div style={{ padding: '4px 12px', background: 'rgba(34, 197, 94, 0.1)', borderRadius: '20px', color: '#22c55e', fontSize: '0.85rem', fontWeight: 600 }}>
                LUNAS
              </motion.div>
            )}
            {paymentStatus === 'failed' && (
              <motion.div style={{ padding: '4px 12px', background: 'rgba(239, 68, 68, 0.1)', borderRadius: '20px', color: '#ef4444', fontSize: '0.85rem', fontWeight: 600 }}>
                KEDALUWARSA
              </motion.div>
            )}
          </div>

          {paymentStatus === 'failed' && (
            <p style={{ color: 'var(--text-dim)', fontSize: '0.85rem' }}>
              Batas waktu pembayaran telah habis. Anda dapat membuat tagihan baru untuk melanjutkan pendaftaran.
            </p>
          )}

          {canCreateInvoice && (
            <button onClick={handleGenerateInvoice} disabled={isLoading} className="btn-primary" style={{ width: '100%', justifyContent: 'center' }}>
              {isLoading ? 'Memproses...' : (
                <>
                  <RefreshCw size={18} />
                  {paymentStatus === 'failed' ? 'Buat Tagihan Baru' : 'Generate Tagihan Pembayaran'}
                </>
              )}
            </button>
          )}
          
          {paymentStatus === 'pending' && (
            <>
              <button onClick={handleOpenPayment} disabled={isLoading || !payment?.payment_url} className="btn-primary" style={{ width: '100%', justifyContent: 'center', background: '#eab308', color: '#000', boxShadow: '0 0 15px rgba(234, 179, 8, 0.3)' }}>
                Buka Pembayaran Midtrans
              </button>
              <button type="button" onClick={handleSimulasiMidtrans} disabled={isLoading} className="btn-outline" style={{ width: '100%', justifyContent: 'center', fontSize: '0.85rem' }}>
                {isLoading ? 'Memproses...' : 'Simulasi: Bayar Berhasil (dev)'}
              </button>
              <button type="button" onClick={handleSimulasiExpire} disabled={isLoading} className="btn-outline" style={{ width: '100%', justifyContent: 'center', fontSize: '0.85rem', borderColor: 'rgba(239,68,68,0.5)', color: '#ef4444' }}>
                Simulasi: Lewat Batas Bayar (dev)
              </button>
            </>
          )}
        </motion.div>

        <motion.div variants={itemVariants} className="glass-card" style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
            <div style={{ padding: '12px', background: paymentStatus === 'success' ? 'rgba(34, 211, 238, 0.1)' : 'rgba(239, 68, 68, 0.1)', borderRadius: '12px' }}>
              <MapPin color={paymentStatus === 'success' ? 'var(--primary)' : '#ef4444'} size={24} />
            </div>
            <div>
              <h3 style={{ fontSize: '1.1rem' }}>Plotting Lokasi</h3>
              <p style={{ color: 'var(--text-dim)', fontSize: '0.9rem' }}>Penempatan Desa</p>
            </div>
          </div>
          
          {!profile?.lokasi ? (
            <div style={{ padding: '1.5rem', background: 'rgba(255,255,255,0.03)', borderRadius: '12px', textAlign: 'center', display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '0.5rem' }}>
              <AlertCircle color="var(--text-dim)" size={28} />
              <p style={{ color: 'var(--text-dim)', fontSize: '0.9rem' }}>
                {paymentStatus === 'success'
                  ? 'Lokasi akan diumumkan setelah proses plotting admin selesai.'
                  : 'Lokasi dan jadwal keberangkatan akan diumumkan setelah pembayaran diverifikasi.'}
              </p>
            </div>
          ) : (
            <div style={{ padding: '1.5rem', background: 'rgba(34, 211, 238, 0.05)', border: '1px solid var(--border)', borderRadius: '12px', display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
              <p style={{ fontSize: '0.85rem', color: 'var(--primary)' }}>Selamat! Anda ditempatkan di:</p>
              <h3 style={{ fontSize: '1.3rem', fontWeight: 700 }}>{profile.lokasi.nama_desa}</h3>
              <p style={{ color: 'var(--text-dim)', fontSize: '0.9rem' }}>Kecamatan {profile.lokasi.kecamatan}</p>
            </div>
          )}
        </motion.div>

        <motion.div variants={itemVariants} className="glass-card" style={{ gridColumn: '1 / -1', marginTop: '1rem' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '1rem', marginBottom: '1.5rem' }}>
            <Bell color="var(--primary)" size={20} />
            <h3 style={{ fontSize: '1.2rem' }}>Papan Informasi Penting</h3>
          </div>
          
          <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
            {notifikasi.length === 0 ? (
              <p style={{ color: 'var(--text-dim)', fontSize: '0.9rem' }}>Belum ada notifikasi.</p>
            ) : (
              notifikasi.map((n) => (
                <div
                  key={n.id}
                  style={{
                    padding: '1rem',
                    borderLeft: `4px solid ${n.is_read ? 'var(--border)' : 'var(--primary)'}`,
                    background: 'rgba(34, 211, 238, 0.05)',
                    borderRadius: '0 12px 12px 0',
                    opacity: n.is_read ? 0.7 : 1,
                  }}
                >
                  <h4 style={{ fontSize: '1rem', marginBottom: '0.3rem' }}>{n.title}</h4>
                  <p style={{ color: 'var(--text-dim)', fontSize: '0.9rem' }}>{n.message}</p>
                </div>
              ))
            )}
          </div>
        </motion.div>

      </motion.div>
    </motion.div>
  );
};

export default Dashboard;
