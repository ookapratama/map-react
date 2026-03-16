import React, { useState } from "react";
import useLocalStorage from "./hooks/useLocalStorage";
import { INITIAL_SPOTS } from "./constants/spots";
import ParkingMap from "./components/ParkingMap";
import BookingForm from "./components/BookingForm";
import BookingDetails from "./components/BookingDetails";
import {
  Car,
  LayoutDashboard,
  History,
  CheckCircle2,
  Menu,
} from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";

function App() {
  const [spots] = useState(INITIAL_SPOTS);
  const [bookings, setBookings] = useLocalStorage("parking_bookings", []);
  const [selectedSpot, setSelectedSpot] = useState(null);
  const [showForm, setShowForm] = useState(false);
  const [showNotification, setShowNotification] = useState(false);

  const activeBookings = bookings.filter((b) => b.status === "active");
  const availableCount = spots.length - activeBookings.length;

  const handleSpotClick = (spot) => {
    const activeBooking = activeBookings.find((b) => b.spotId === spot.id);
    if (activeBooking) {
      setSelectedSpot({ spot, booking: activeBooking });
    } else {
      setSelectedSpot({ spot, booking: null });
      setShowForm(true);
    }
  };

  const handleBook = (newBooking) => {
    const booking = {
      ...newBooking,
      id: Date.now().toString(),
    };
    setBookings([...bookings, booking]);
    setShowForm(false);
    setShowNotification(true);
    setTimeout(() => setShowNotification(false), 3000);
  };

  const handleEndSession = (bookingId) => {
    const updatedBookings = bookings.map((b) =>
      b.id === bookingId
        ? { ...b, status: "completed", endTime: new Date().toISOString() }
        : b,
    );
    setBookings(updatedBookings);
    setSelectedSpot(null);
  };

  return (
    <div className="flex min-h-screen bg-slate-50 font-sans">
      {/* Sidebar - Desktop Only */}
      <aside className="hidden lg:flex w-64 bg-white border-r border-slate-200 p-8 flex-col gap-10">
        <div className="flex items-center gap-3">
          <div className="bg-primary/10 p-2 rounded-lg">
            <Car size={32} className="text-primary" />
          </div>
          <span className="font-heading text-2xl font-extrabold tracking-tight">
            ParkSpot
          </span>
        </div>
        <nav className="flex flex-col gap-2">
          <button className="flex items-center gap-3 px-4 py-3 rounded-xl bg-primary/5 text-primary font-bold transition-all">
            <LayoutDashboard size={20} /> Dashboard
          </button>
          <button className="flex items-center gap-3 px-4 py-3 rounded-xl text-slate-500 hover:bg-slate-50 hover:text-slate-900 font-semibold transition-all">
            <History size={20} /> Riwayat
          </button>
        </nav>
      </aside>

      <main className="flex-1 overflow-y-auto">
        <header className="px-6 lg:px-10 py-8 lg:py-10 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-6 max-w-7xl mx-auto w-full">
          <div>
            <h1 className="text-3xl lg:text-4xl">Pengelolaan Parkiran</h1>
            <p className="text-slate-500 mt-1">
              Pantau dan kelola ketersediaan tempat parkir secara real-time.
            </p>
          </div>
          <div className="flex gap-4 w-full sm:w-auto">
            <div className="bg-white p-4 lg:px-6 rounded-xl shadow-sm border border-slate-200 flex flex-col min-w-[120px] flex-1 sm:flex-none">
              <span className="text-[10px] lg:text-xs font-bold text-slate-400 uppercase tracking-wider">
                Tersedia
              </span>
              <span className="text-2xl lg:text-3xl font-heading font-black text-success leading-tight">
                {availableCount}
              </span>
            </div>
            <div className="bg-white p-4 lg:px-6 rounded-xl shadow-sm border border-slate-200 flex flex-col min-w-[120px] flex-1 sm:flex-none">
              <span className="text-[10px] lg:text-xs font-bold text-slate-400 uppercase tracking-wider">
                Terisi
              </span>
              <span className="text-2xl lg:text-3xl font-heading font-black text-danger leading-tight">
                {activeBookings.length}
              </span>
            </div>
          </div>
        </header>

        <section className="px-6 lg:px-10 pb-10 grid grid-cols-1 xl:grid-cols-[1fr_380px] gap-8 max-w-7xl mx-auto w-full">
          <div className="space-y-6">
            <div className="flex justify-between items-center">
              <h2 className="text-xl">Denah Parkiran</h2>
              <div className="flex gap-4 text-sm font-medium">
                <div className="flex items-center gap-2 text-slate-500">
                  <div className="w-3 h-3 rounded-full bg-success"></div>{" "}
                  Tersedia
                </div>
                <div className="flex items-center gap-2 text-slate-500">
                  <div className="w-3 h-3 rounded-full bg-danger"></div> Terisi
                </div>
              </div>
            </div>

            <div className="bg-white p-4 lg:p-8 rounded-2xl shadow-sm border border-slate-200 min-h-[450px] flex items-center justify-center overflow-auto">
              <ParkingMap
                spots={spots}
                bookings={bookings}
                onSpotClick={handleSpotClick}
              />
            </div>
          </div>

          <aside className="space-y-6">
            {selectedSpot && selectedSpot.booking ? (
              <div className="sticky top-10 animate-in fade-in slide-in-from-right-4 duration-300">
                <h2 className="text-xl mb-4">Rincian Pemesanan</h2>
                <BookingDetails
                  booking={selectedSpot.booking}
                  spot={selectedSpot.spot}
                  onEndSession={handleEndSession}
                />
              </div>
            ) : (
              <div className="bg-white rounded-2xl p-10 border border-slate-200 border-dashed flex flex-col items-center justify-center text-center gap-4 text-slate-400 sticky top-10">
                <div className="w-16 h-16 bg-slate-50 rounded-full flex items-center justify-center">
                  <Car size={32} />
                </div>
                <div>
                  <h3 className="text-slate-900 text-lg">
                    Pilih Tempat Parkir
                  </h3>
                  <p className="text-sm max-w-[200px] mx-auto">
                    Klik pada tempat parkir di denah untuk melihat rincian.
                  </p>
                </div>
              </div>
            )}
          </aside>
        </section>
      </main>

      <AnimatePresence>
        {showForm && selectedSpot && (
          <BookingForm
            spot={selectedSpot.spot}
            onClose={() => setShowForm(false)}
            onBook={handleBook}
          />
        )}
      </AnimatePresence>

      <AnimatePresence>
        {showNotification && (
          <motion.div
            initial={{ opacity: 0, x: 100 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: 100 }}
            className="fixed bottom-8 right-8 bg-white p-5 rounded-2xl shadow-2xl border-l-4 border-success flex items-center gap-4 z-[100] max-w-sm"
          >
            <div className="bg-success/10 p-2 rounded-full text-success">
              <CheckCircle2 size={24} />
            </div>
            <div>
              <h4 className="font-bold text-slate-900">Pemesanan Berhasil</h4>
              <p className="text-sm text-slate-500">
                Tempat parkir telah berhasil dipesan.
              </p>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}

export default App;
