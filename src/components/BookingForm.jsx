import React, { useState } from "react";
import { X } from "lucide-react";
import { motion } from "framer-motion";

const BookingForm = ({ spot, onClose, onBook }) => {
  const [formData, setFormData] = useState({
    name: "",
    vehicleNumber: "",
    duration: 60, // minutes
  });

  const handleSubmit = (e) => {
    e.preventDefault();
    onBook({
      ...formData,
      spotId: spot.id,
      startTime: new Date().toISOString(),
      status: "active",
    });
  };

  return (
    <div className="fixed inset-0 bg-slate-900/40 backdrop-blur-sm z-[1000] flex items-center justify-center p-4">
      <motion.div
        initial={{ opacity: 0, scale: 0.95, y: 20 }}
        animate={{ opacity: 1, scale: 1, y: 0 }}
        exit={{ opacity: 0, scale: 0.95, y: 20 }}
        className="w-full max-w-md bg-white rounded-3xl shadow-2xl overflow-hidden"
      >
        <div className="p-6 lg:p-8 space-y-6">
          <div className="flex justify-between items-center">
            <h3 className="text-2xl font-heading font-bold text-slate-900">
              Pesan <span className="text-primary">{spot.label}</span>
            </h3>
            <button
              onClick={onClose}
              className="p-2 hover:bg-slate-100 rounded-full text-slate-400 hover:text-slate-900 transition-colors"
            >
              <X size={24} />
            </button>
          </div>

          <form onSubmit={handleSubmit} className="space-y-5">
            <div className="space-y-2">
              <label className="text-sm font-bold text-slate-700 uppercase tracking-tight">
                Nama Pengguna
              </label>
              <input
                type="text"
                required
                placeholder="Masukkan nama lengkap"
                className="w-full px-4 py-3 rounded-xl border border-slate-200 focus:border-primary focus:ring-4 focus:ring-primary/10 outline-none transition-all placeholder:text-slate-300"
                value={formData.name}
                onChange={(e) =>
                  setFormData({ ...formData, name: e.target.value })
                }
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-bold text-slate-700 uppercase tracking-tight">
                Nomor Kendaraan
              </label>
              <input
                type="text"
                required
                placeholder="Contoh: B 1234 ABC"
                className="w-full px-4 py-3 rounded-xl border border-slate-200 focus:border-primary focus:ring-4 focus:ring-primary/10 outline-none transition-all placeholder:text-slate-300"
                value={formData.vehicleNumber}
                onChange={(e) =>
                  setFormData({ ...formData, vehicleNumber: e.target.value })
                }
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-bold text-slate-700 uppercase tracking-tight">
                Durasi Parkir
              </label>
              <select
                className="w-full px-4 py-3 rounded-xl border border-slate-200 focus:border-primary focus:ring-4 focus:ring-primary/10 outline-none transition-all appearance-none bg-no-repeat bg-[right_1rem_center]"
                value={formData.duration}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    duration: parseInt(e.target.value),
                  })
                }
                style={{
                  backgroundImage: `url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 24 24' stroke='%2364748b' stroke-width='2'%3E%3Cpath stroke-linecap='round' stroke-linejoin='round' d='M19 9l-7 7-7-7'%3E%3C/path%3E%3C/svg%3E")`,
                  backgroundSize: "1.25rem",
                }}
              >
                <option value={30}>30 Menit</option>
                <option value={60}>1 Jam</option>
                <option value={120}>2 Jam</option>
                <option value={180}>3 Jam</option>
                <option value={300}>5 Jam</option>
              </select>
            </div>

            <div className="flex gap-4 pt-4">
              <button
                type="button"
                className="flex-1 px-6 py-3.5 rounded-xl border border-slate-200 font-bold text-slate-600 hover:bg-slate-50 transition-colors"
                onClick={onClose}
              >
                Batal
              </button>
              <button
                type="submit"
                className="flex-1 px-6 py-3.5 rounded-xl bg-primary text-white font-bold hover:bg-primary-hover shadow-lg shadow-primary/30 transition-all hover:-translate-y-0.5"
              >
                Konfirmasi
              </button>
            </div>
          </form>
        </div>
      </motion.div>
    </div>
  );
};

export default BookingForm;
