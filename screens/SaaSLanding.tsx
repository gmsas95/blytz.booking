import React, { useState, useEffect } from 'react';
import { ShieldCheck, ArrowRight, Zap, Globe, LayoutGrid, CheckCircle2, TrendingUp } from 'lucide-react';
import { api, Business } from '../api';

interface SaaSLandingProps {
  onSelectBusiness: (business: Business) => void;
  onOperatorLogin: () => void;
}

export const SaaSLanding: React.FC<SaaSLandingProps> = ({ onSelectBusiness, onOperatorLogin }) => {
  const [businesses, setBusinesses] = useState<Business[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchBusinesses = async () => {
      try {
        setLoading(true);
        const data = await api.getBusinesses();
        setBusinesses(data);
      } catch (err) {
        console.error('Failed to fetch businesses:', err);
        setError('Failed to load businesses. Please try again later.');
      } finally {
        setLoading(false);
      }
    };

    fetchBusinesses();
  }, []);

  return (
    <div className="min-h-screen bg-zinc-950 text-zinc-50 font-sans selection:bg-primary-500 selection:text-black">
      {/* Background Grid Texture */}
      <div className="fixed inset-0 z-0 pointer-events-none opacity-[0.07]"
           style={{
             backgroundImage: 'linear-gradient(#fff 1px, transparent 1px), linear-gradient(90deg, #fff 1px, transparent 1px)',
             backgroundSize: '50px 50px'
           }}
      />

      {/* Header */}
      <header className="relative z-10 border-b border-zinc-800 bg-zinc-950/80 backdrop-blur-md">
        <div className="max-w-7xl mx-auto px-6 h-20 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="bg-primary-500 h-8 w-8 flex items-center justify-center rounded-sm text-black font-black text-xl">
              B
            </div>
            <span className="text-xl font-bold tracking-tighter">BLYTZ<span className="text-zinc-500">.CLOUD</span></span>
          </div>
          <div className="flex items-center gap-6">
            <button
                onClick={onOperatorLogin}
                className="hidden sm:flex text-sm font-medium text-zinc-400 hover:text-white transition-colors"
            >
              Operator Login
            </button>
            <button className="bg-zinc-100 hover:bg-white text-black text-sm font-bold py-2.5 px-5 transition-transform hover:-translate-y-0.5 active:translate-y-0 rounded-sm">
              GET STARTED
            </button>
          </div>
        </div>
      </header>

      {/* Hero Section */}
      <section className="relative z-10 pt-24 pb-20 px-6 border-b border-zinc-800">
        <div className="max-w-7xl mx-auto text-center">
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full border border-zinc-800 bg-zinc-900/50 text-primary-400 text-xs font-mono mb-8">
            <div className="w-2 h-2 rounded-full bg-primary-500 animate-pulse" />
            V2.0 LIVE: MULTI-TENANT ENGINE
          </div>

          <h1 className="text-6xl sm:text-8xl md:text-9xl font-black tracking-tighter text-white leading-[0.9] mb-8">
            NO DEPOSIT.<br />
            <span className="text-zinc-800 stroke-text">NO BOOKING.</span>
          </h1>

          <p className="text-lg sm:text-xl text-zinc-400 max-w-2xl mx-auto mb-10 leading-relaxed font-light">
            Stop chasing invoices. The cloud-based booking management solution for freelancers that forces upfront payment.
            If they don't pay, they don't get the slot. Simple.
          </p>

          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
             <button onClick={() => document.getElementById('demos')?.scrollIntoView({behavior: 'smooth'})} className="h-14 px-8 bg-primary-500 hover:bg-primary-400 text-black text-lg font-bold rounded-sm w-full sm:w-auto transition-all">
               TRY THE DEMO
             </button>
             <button className="h-14 px-8 border border-zinc-700 hover:border-zinc-500 hover:bg-zinc-900 text-white text-lg font-medium rounded-sm w-full sm:w-auto transition-all">
               VIEW DOCS
             </button>
          </div>
        </div>
      </section>

      {/* Marquee / Stats Strip */}
      <div className="relative z-10 border-b border-zinc-800 bg-zinc-900/50 overflow-hidden py-4">
        <div className="flex items-center justify-center gap-8 sm:gap-16 text-zinc-500 font-mono text-xs sm:text-sm uppercase tracking-widest opacity-70">
            <span className="flex items-center gap-2"><CheckCircle2 className="h-4 w-4 text-primary-500" /> No Ghosting</span>
            <span className="hidden sm:flex items-center gap-2"><CheckCircle2 className="h-4 w-4 text-primary-500" /> Stripe Connect</span>
            <span className="flex items-center gap-2"><CheckCircle2 className="h-4 w-4 text-primary-500" /> Instant Cashflow</span>
            <span className="hidden sm:flex items-center gap-2"><CheckCircle2 className="h-4 w-4 text-primary-500" /> Zero Fluff</span>
        </div>
      </div>

      {/* Interactive Demo Section */}
      <section id="demos" className="relative z-10 py-24 px-6 max-w-7xl mx-auto">
        <div className="flex flex-col md:flex-row md:items-end justify-between mb-12 gap-6">
          <div>
            <h2 className="text-4xl sm:text-5xl font-bold tracking-tighter text-white mb-4">
              CHOOSE YOUR FIGHTER
            </h2>
            <p className="text-zinc-400 max-w-md">
              See how Blytz.Cloud adapts to any service vertical. Select a preset below to enter the booking flow.
            </p>
          </div>
          <div className="text-right hidden md:block">
            <div className="text-zinc-500 font-mono text-sm">{businesses.length} VERTICALS AVAILABLE</div>
          </div>
        </div>

        {loading ? (
          <div className="text-center py-20 text-zinc-500">Loading businesses...</div>
        ) : error ? (
          <div className="text-center py-20 text-red-500">{error}</div>
        ) : businesses.length === 0 ? (
          <div className="text-center py-20 text-zinc-500">No businesses available yet.</div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-0 border border-zinc-800 bg-zinc-900">
            {businesses.map((biz, idx) => (
              <div
                key={biz.id}
                onClick={() => onSelectBusiness(biz)}
                className={`
                  group relative p-8 cursor-pointer transition-all duration-300 hover:bg-zinc-800
                  ${idx !== businesses.length - 1 ? 'border-b md:border-b-0 md:border-r border-zinc-800' : ''}
                `}
              >
                <div className="absolute top-4 right-4 opacity-0 group-hover:opacity-100 transition-opacity">
                  <ArrowRight className="text-primary-500 -rotate-45 group-hover:rotate-0 transition-transform duration-300" />
                </div>

                <div className="mb-8">
                    <div className={`
                      w-12 h-12 flex items-center justify-center rounded-sm text-lg font-bold mb-4
                      ${biz.vertical === 'Automotive' ? 'bg-blue-600 text-white' :
                        biz.vertical === 'Wellness' ? 'bg-emerald-500 text-black' :
                        biz.vertical === 'Creative' ? 'bg-violet-600 text-white' : 'bg-zinc-100 text-black'
                      }
                    `}>
                      {biz.name.charAt(0)}
                    </div>
                    <div className="inline-block px-2 py-1 bg-zinc-950 text-zinc-500 border border-zinc-800 text-[10px] font-mono uppercase tracking-wider mb-2">
                      {biz.vertical}
                    </div>
                    <h3 className="text-xl font-bold text-white group-hover:text-primary-500 transition-colors">
                      {biz.name}
                    </h3>
                </div>

                <p className="text-zinc-400 text-sm leading-relaxed mb-8 h-10">
                  {biz.description}
                </p>

                <div className="flex items-center text-xs font-bold tracking-widest uppercase text-zinc-500 group-hover:text-white transition-colors">
                   <span>Enter Booking Flow</span>
                   <div className="h-[1px] bg-zinc-700 w-8 ml-3 group-hover:w-16 group-hover:bg-primary-500 transition-all"></div>
                </div>
              </div>
            ))}
          </div>
        )}
      </section>

      {/* Grid Features */}
      <section className="border-t border-zinc-800 bg-zinc-950 relative z-10">
        <div className="max-w-7xl mx-auto grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 divide-y md:divide-y-0 md:divide-x divide-zinc-800 border-x border-zinc-800">
           {/* Box 1 */}
           <div className="p-10 hover:bg-zinc-900 transition-colors">
              <Zap className="h-8 w-8 text-white mb-6" />
              <h3 className="text-lg font-bold text-white mb-2">Lightning Fast</h3>
              <p className="text-sm text-zinc-500">Built on Next.js 14. Zero bloat. Loads instantly on 3G connections.</p>
           </div>
           {/* Box 2 */}
           <div className="p-10 hover:bg-zinc-900 transition-colors">
              <ShieldCheck className="h-8 w-8 text-white mb-6" />
              <h3 className="text-lg font-bold text-white mb-2">Fraud Proof</h3>
              <p className="text-sm text-zinc-500">Stripe integration handles 3D Secure and prevents card testing attacks.</p>
           </div>
           {/* Box 3 */}
           <div className="p-10 hover:bg-zinc-900 transition-colors">
              <LayoutGrid className="h-8 w-8 text-white mb-6" />
              <h3 className="text-lg font-bold text-white mb-2">Infinite Scales</h3>
              <p className="text-sm text-zinc-500">One dashboard. Unlimited locations. Centralized revenue tracking.</p>
           </div>
           {/* Box 4 */}
           <div className="p-10 bg-zinc-900/50 hover:bg-zinc-900 transition-colors flex flex-col justify-center">
              <div className="text-4xl font-black text-white mb-2">$4.2M</div>
              <p className="text-sm text-zinc-500 font-mono">PROCESSED IN DEPOSITS</p>
              <div className="h-1 w-full bg-zinc-800 mt-6 overflow-hidden">
                <div className="h-full bg-primary-500 w-3/4"></div>
              </div>
           </div>
        </div>
      </section>

      <footer className="py-12 px-6 border-t border-zinc-800 bg-zinc-950 text-center relative z-10">
        <p className="text-zinc-600 text-sm font-mono">
          © 2024 BLYTZ.CLOUD INC. SHIP FAST OR DIE TRYING.
        </p>
      </footer>

      {/* CSS Helper for stroke text effect since tailwind doesn't have it native */}
      <style>{`
        .stroke-text {
          -webkit-text-stroke: 1px #3f3f46;
          color: transparent;
        }
      `}</style>
    </div>
  );
};
