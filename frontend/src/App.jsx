import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { 
  Users, 
  MapPin, 
  ShieldCheck, 
  ArrowRight, 
  LayoutGrid, 
  Bell, 
  CreditCard,
  LogOut
} from 'lucide-react';
import LoginModal from './components/ui/LoginModal';
import PanduanModal from './components/ui/PanduanModal';
import Dashboard from './components/layout/Dashboard';

const App = () => {
  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false);
  const [authModalView, setAuthModalView] = useState('login'); // 'login' | 'register'
  const [isPanduanOpen, setIsPanduanOpen] = useState(false);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [userNim, setUserNim] = useState('');

  // Cek token saat komponen dimuat (Sesi 30 hari)
  useEffect(() => {
    const token = localStorage.getItem('kkn_token');
    const nim = localStorage.getItem('kkn_nim');
    if (token) {
      setIsLoggedIn(true);
      setUserNim(nim);
    }
  }, []);

  const handleLogout = () => {
    localStorage.removeItem('kkn_token');
    localStorage.removeItem('kkn_nim');
    localStorage.removeItem('kkn_peserta_id');
    setIsLoggedIn(false);
  };

  const handleLoginSuccess = () => {
    setIsLoggedIn(true);
    setUserNim(localStorage.getItem('kkn_nim'));
    setIsLoginModalOpen(false);
  };
  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.2
      }
    }
  };

  const itemVariants = {
    hidden: { y: 20, opacity: 0 },
    visible: { y: 0, opacity: 1 }
  };

  return (
    <div className="min-h-screen relative overflow-hidden">
      {/* Decorative Glows */}
      <div className="hero-glow" style={{ top: '-10%', left: '-10%' }} />
      <div className="hero-glow" style={{ bottom: '-10%', right: '-10%' }} />

      {/* Navbar */}
      <nav style={{ 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center', 
        padding: '1.5rem 5%',
        position: 'relative',
        zIndex: 10
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '0.8rem' }}>
          <div style={{ 
            background: 'var(--primary)', 
            padding: '8px', 
            borderRadius: '12px',
            boxShadow: '0 0 15px var(--primary)'
          }}>
            <LayoutGrid size={24} color="#020617" />
          </div>
          <span style={{ fontWeight: 700, fontSize: '1.4rem', color: '#FFF' }}>
            KKN<span style={{ color: 'var(--primary)' }}>Hub</span>
          </span>
        </div>
        <div style={{ display: 'flex', gap: '2rem', alignItems: 'center' }}>
          <a href="#" style={{ color: 'var(--text-dim)', textDecoration: 'none', fontSize: '0.9rem' }}>Beranda</a>
          <a href="#" style={{ color: 'var(--text-dim)', textDecoration: 'none', fontSize: '0.9rem' }}>Informasi</a>
          
          {isLoggedIn ? (
            <button onClick={handleLogout} className="btn-outline" style={{ padding: '0.6rem 1.2rem', fontSize: '0.9rem', borderColor: '#ef4444', color: '#ef4444' }}>
              Keluar ({userNim}) <LogOut size={16} />
            </button>
          ) : (
            <button onClick={() => { setAuthModalView('login'); setIsLoginModalOpen(true); }} className="btn-outline" style={{ padding: '0.6rem 1.2rem', fontSize: '0.9rem' }}>
              Login Portal
            </button>
          )}
        </div>
      </nav>

      {/* Konten Berubah Tergantung Status Login */}
      {isLoggedIn ? (
        <Dashboard userNim={userNim} onLogout={handleLogout} />
      ) : (
        <main style={{ padding: '4rem 5%', position: 'relative', zIndex: 1 }}>
        <div style={{ maxWidth: '1200px', margin: '0 auto' }}>
          <div style={{ display: 'grid', gridTemplateColumns: '1.2fr 0.8fr', gap: '4rem', alignItems: 'center' }}>
            <motion.div 
              initial={{ x: -50, opacity: 0 }}
              animate={{ x: 0, opacity: 1 }}
              transition={{ duration: 0.8 }}
            >
              <span style={{ 
                color: 'var(--primary)', 
                fontWeight: 600, 
                fontSize: '0.9rem', 
                letterSpacing: '2px', 
                textTransform: 'uppercase',
                display: 'block',
                marginBottom: '1rem'
              }}>
                Platform Manajemen KKN Terintegrasi
              </span>
              <h1 className="gradient-text" style={{ fontSize: '4rem', fontWeight: 800, lineHeight: 1.1, marginBottom: '1.5rem' }}>
                Membangun Negeri, <br /> Dimulai dari Desa.
              </h1>
              <p style={{ color: 'var(--text-dim)', fontSize: '1.2rem', maxWidth: '500px', marginBottom: '2.5rem' }}>
                Sistem pendaftaran dan pengelolaan KKN yang transparan, efisien, dan modern untuk seluruh mahasiswa.
              </p>
              <div style={{ display: 'flex', gap: '1.5rem' }}>
                <button onClick={() => { setAuthModalView('register'); setIsLoginModalOpen(true); }} className="btn-primary">
                  Mulai Pendaftaran <ArrowRight size={20} />
                </button>
                <button onClick={() => setIsPanduanOpen(true)} className="btn-outline">
                  Panduan Sistem
                </button>
              </div>
            </motion.div>

            <motion.div 
              initial={{ scale: 0.8, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              transition={{ duration: 1 }}
              className="glass-card"
              style={{ position: 'relative' }}
            >
              <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <h3 style={{ fontSize: '1.2rem' }}>Informasi Terkini</h3>
                  <Bell size={20} color="var(--primary)" />
                </div>
                <div style={{ padding: '1rem', background: 'rgba(255,255,255,0.03)', borderRadius: '16px', border: '1px solid var(--border)' }}>
                  <p style={{ fontSize: '0.8rem', color: 'var(--text-dim)' }}>Status Pendaftaran</p>
                  <p style={{ fontWeight: 600 }}>Menunggu Verifikasi</p>
                </div>
                <div style={{ padding: '1rem', background: 'rgba(34, 211, 238, 0.05)', borderRadius: '16px', border: '1px solid var(--primary)' }}>
                  <p style={{ fontSize: '0.8rem', color: 'var(--primary)' }}>Sisa Kuota Lokasi</p>
                  <p style={{ fontWeight: 600 }}>124 Mahasiswa Lagi</p>
                </div>
              </div>
            </motion.div>
          </div>

          {/* Features Grid */}
          <motion.div 
            variants={containerVariants}
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true }}
            style={{ 
              marginTop: '6rem',
              display: 'grid', 
              gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))', 
              gap: '2rem' 
            }}
          >
            {[
              { icon: Users, title: "Pendaftaran Cepat", desc: "Proses pendaftaran kolektif maupun mandiri yang simpel." },
              { icon: MapPin, title: "Plotting Cerdas", desc: "Penempatan lokasi berdasarkan minat dan kebutuhan desa." },
              { icon: ShieldCheck, title: "Verifikasi Berkas", desc: "Keamanan data dan validasi dokumen secara digital." },
              { icon: CreditCard, title: "Pembayaran Mudah", desc: "Integrasi payment gateway untuk biaya administrasi." }
            ].map((f, i) => (
              <motion.div key={i} variants={itemVariants} className="glass-card">
                <div style={{ 
                  width: '50px', 
                  height: '50px', 
                  background: 'rgba(34, 211, 238, 0.1)', 
                  borderRadius: '12px', 
                  display: 'flex', 
                  alignItems: 'center', 
                  justifyContent: 'center',
                  marginBottom: '1.5rem'
                }}>
                  <f.icon color="var(--primary)" size={24} />
                </div>
                <h4 style={{ fontSize: '1.3rem', marginBottom: '0.8rem' }}>{f.title}</h4>
                <p style={{ color: 'var(--text-dim)', fontSize: '0.95rem' }}>{f.desc}</p>
              </motion.div>
            ))}
          </motion.div>
        </div>
      </main>
      )}

      {/* Modal Autentikasi (Akan dirender jika state isLoginModalOpen = true) */}
      <AnimatePresence>
        {isLoginModalOpen && (
          <LoginModal 
            isOpen={isLoginModalOpen} 
            onClose={() => setIsLoginModalOpen(false)} 
            onLoginSuccess={handleLoginSuccess}
            defaultView={authModalView}
          />
        )}
      </AnimatePresence>

      {/* Modal Panduan */}
      <AnimatePresence>
        {isPanduanOpen && (
          <PanduanModal 
            isOpen={isPanduanOpen} 
            onClose={() => setIsPanduanOpen(false)} 
          />
        )}
      </AnimatePresence>

      {/* Background Shapes */}
      <div style={{
        position: 'absolute',
        top: '20%',
        right: '-5%',
        width: '500px',
        height: '500px',
        borderRadius: '50%',
        background: 'rgba(34, 211, 238, 0.02)',
        border: '1px solid rgba(34, 211, 238, 0.05)',
        zIndex: -1
      }} />
    </div>
  );
};

export default App;
