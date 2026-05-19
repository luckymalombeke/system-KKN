import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { X, ArrowRight, CheckCircle2, Loader2, UserPlus, LogIn } from 'lucide-react';
import { KKN_API } from '../../services/api';

const inputStyle = {
  width: '100%',
  padding: '1rem',
  borderRadius: '12px',
  background: 'rgba(255,255,255,0.05)',
  border: '1px solid var(--border)',
  color: 'white',
  fontSize: '1rem',
  outline: 'none',
};

const LoginModal = ({ isOpen, onClose, onLoginSuccess, defaultView = 'login' }) => {
  const [view, setView] = useState(defaultView); // 'login' | 'register'
  const [loginStep, setLoginStep] = useState('input-nim'); // 'input-nim' | 'input-otp'
  const [nim, setNim] = useState('');
  const [otp, setOtp] = useState('');
  const [registerForm, setRegisterForm] = useState({
    nama: '',
    email: '',
    nim: '',
    prodi: '',
  });
  const [successMessage, setSuccessMessage] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    if (isOpen) {
      setView(defaultView);
      setLoginStep('input-nim');
      setError('');
      setSuccessMessage('');
    }
  }, [isOpen, defaultView]);

  const handleRequestOtp = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');
    setSuccessMessage('');

    try {
      await KKN_API.requestOtp(nim.trim());
      setLoginStep('input-otp');
    } catch (err) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  const handleVerifyOtp = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');

    try {
      const res = await KKN_API.verifyOtp(nim.trim(), otp.trim());
      localStorage.setItem('kkn_token', res.token);
      localStorage.setItem('kkn_nim', nim.trim());
      localStorage.setItem('kkn_peserta_id', res.peserta_id);
      onLoginSuccess();
    } catch (err) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  const handleRegister = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');
    setSuccessMessage('');

    try {
      const payload = {
        nama: registerForm.nama.trim(),
        email: registerForm.email.trim(),
        nim: registerForm.nim.trim(),
        prodi: registerForm.prodi.trim(),
      };
      await KKN_API.registerPeserta(payload);
      setNim(payload.nim);
      setSuccessMessage('Pendaftaran berhasil! Silakan masuk dengan OTP.');
      setView('login');
      setLoginStep('input-nim');
    } catch (err) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  const switchView = (nextView) => {
    setView(nextView);
    setError('');
    setSuccessMessage('');
    setLoginStep('input-nim');
  };

  if (!isOpen) return null;

  return (
    <motion.div
      style={{
        position: 'fixed', top: 0, left: 0, right: 0, bottom: 0,
        background: 'rgba(2, 6, 23, 0.8)',
        backdropFilter: 'blur(8px)',
        display: 'flex', justifyContent: 'center', alignItems: 'center',
        zIndex: 100,
      }}
    >
      <motion.div
        initial={{ scale: 0.9, opacity: 0, y: 20 }}
        animate={{ scale: 1, opacity: 1, y: 0 }}
        exit={{ scale: 0.9, opacity: 0, y: 20 }}
        className="glass-card"
        style={{ width: '100%', maxWidth: '440px', position: 'relative', maxHeight: '90vh', overflowY: 'auto' }}
      >
        <button
          type="button"
          onClick={onClose}
          style={{ position: 'absolute', top: '15px', right: '15px', background: 'transparent', border: 'none', color: 'var(--text-dim)', cursor: 'pointer' }}
        >
          <X size={24} />
        </button>

        <motion.div
          style={{
            display: 'flex',
            gap: '0.5rem',
            marginBottom: '1.5rem',
            padding: '4px',
            background: 'rgba(255,255,255,0.05)',
            borderRadius: '12px',
          }}
        >
          <button
            type="button"
            onClick={() => switchView('login')}
            className={view === 'login' ? 'btn-primary' : 'btn-outline'}
            style={{ flex: 1, justifyContent: 'center', padding: '0.6rem', fontSize: '0.9rem' }}
          >
            <LogIn size={16} /> Masuk
          </button>
          <button
            type="button"
            onClick={() => switchView('register')}
            className={view === 'register' ? 'btn-primary' : 'btn-outline'}
            style={{ flex: 1, justifyContent: 'center', padding: '0.6rem', fontSize: '0.9rem' }}
          >
            <UserPlus size={16} /> Daftar KKN
          </button>
        </motion.div>

        <div style={{ textAlign: 'center', marginBottom: '1.5rem' }}>
          <h2 style={{ fontSize: '1.5rem', marginBottom: '0.5rem' }}>
            {view === 'register'
              ? 'Pendaftaran KKN'
              : loginStep === 'input-nim'
                ? 'Masuk ke Portal KKN'
                : 'Verifikasi OTP'}
          </h2>
          <p style={{ color: 'var(--text-dim)', fontSize: '0.9rem' }}>
            {view === 'register'
              ? 'Isi data diri Anda. Setelah itu login dengan NIM dan OTP.'
              : loginStep === 'input-nim'
                ? 'Masukkan NIM yang sudah terdaftar.'
                : `Kode OTP untuk NIM ${nim} telah dikirim ke email terdaftar.`}
          </p>
        </div>

        {successMessage && (
          <motion.div style={{ background: 'rgba(34, 197, 94, 0.1)', color: '#22c55e', padding: '0.8rem', borderRadius: '8px', fontSize: '0.85rem', marginBottom: '1rem', textAlign: 'center' }}>
            {successMessage}
          </motion.div>
        )}

        {error && (
          <motion.div style={{ background: 'rgba(239, 68, 68, 0.1)', color: '#ef4444', padding: '0.8rem', borderRadius: '8px', fontSize: '0.85rem', marginBottom: '1rem', textAlign: 'center' }}>
            {error}
          </motion.div>
        )}

        <AnimatePresence mode="wait">
          {view === 'register' ? (
            <motion.form
              key="register-form"
              initial={{ x: -20, opacity: 0 }}
              animate={{ x: 0, opacity: 1 }}
              exit={{ x: 20, opacity: 0 }}
              onSubmit={handleRegister}
              style={{ display: 'flex', flexDirection: 'column', gap: '0.85rem' }}
            >
              <input
                type="text"
                placeholder="Nama lengkap"
                value={registerForm.nama}
                onChange={(e) => setRegisterForm({ ...registerForm, nama: e.target.value })}
                required
                style={inputStyle}
              />
              <input
                type="email"
                placeholder="Email aktif"
                value={registerForm.email}
                onChange={(e) => setRegisterForm({ ...registerForm, email: e.target.value })}
                required
                style={inputStyle}
              />
              <input
                type="text"
                placeholder="NIM"
                value={registerForm.nim}
                onChange={(e) => setRegisterForm({ ...registerForm, nim: e.target.value })}
                required
                style={inputStyle}
              />
              <input
                type="text"
                placeholder="Program studi"
                value={registerForm.prodi}
                onChange={(e) => setRegisterForm({ ...registerForm, prodi: e.target.value })}
                required
                style={inputStyle}
              />
              <button type="submit" className="btn-primary" style={{ justifyContent: 'center', marginTop: '0.5rem' }} disabled={isLoading}>
                {isLoading ? <Loader2 className="animate-spin" /> : <>Daftar Sekarang <UserPlus size={18} /></>}
              </button>
            </motion.form>
          ) : loginStep === 'input-nim' ? (
            <motion.form
              key="nim-form"
              initial={{ x: -20, opacity: 0 }}
              animate={{ x: 0, opacity: 1 }}
              exit={{ x: 20, opacity: 0 }}
              onSubmit={handleRequestOtp}
              style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}
            >
              <input
                type="text"
                placeholder="Masukkan NIM Anda..."
                value={nim}
                onChange={(e) => setNim(e.target.value)}
                required
                style={inputStyle}
              />
              <button type="submit" className="btn-primary" style={{ justifyContent: 'center' }} disabled={isLoading}>
                {isLoading ? <Loader2 className="animate-spin" /> : <>Kirim Kode OTP <ArrowRight size={18} /></>}
              </button>
              <p style={{ textAlign: 'center', fontSize: '0.85rem', color: 'var(--text-dim)' }}>
                Belum punya akun?{' '}
                <button type="button" onClick={() => switchView('register')} style={{ background: 'none', border: 'none', color: 'var(--primary)', cursor: 'pointer', fontWeight: 600 }}>
                  Daftar di sini
                </button>
              </p>
            </motion.form>
          ) : (
            <motion.form
              key="otp-form"
              initial={{ x: -20, opacity: 0 }}
              animate={{ x: 0, opacity: 1 }}
              exit={{ x: 20, opacity: 0 }}
              onSubmit={handleVerifyOtp}
              style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}
            >
              <input
                type="text"
                placeholder="123456"
                maxLength={6}
                value={otp}
                onChange={(e) => setOtp(e.target.value)}
                required
                style={{ ...inputStyle, border: '1px solid var(--primary)', fontSize: '2rem', letterSpacing: '0.5rem', textAlign: 'center' }}
              />
              <button type="submit" className="btn-primary" style={{ justifyContent: 'center' }} disabled={isLoading}>
                {isLoading ? <Loader2 className="animate-spin" /> : <>Verifikasi & Masuk <CheckCircle2 size={18} /></>}
              </button>
              <button
                type="button"
                onClick={() => setLoginStep('input-nim')}
                style={{ background: 'transparent', border: 'none', color: 'var(--text-dim)', fontSize: '0.85rem', cursor: 'pointer' }}
              >
                Ganti NIM
              </button>
            </motion.form>
          )}
        </AnimatePresence>
      </motion.div>
    </motion.div>
  );
};

export default LoginModal;
