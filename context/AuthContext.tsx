import React, { createContext, useContext, useEffect, useMemo, useState } from 'react';
import { api, CurrentUserResponse, Membership, User } from '../api';

interface AuthContextType {
  isAuthenticated: boolean;
  isLoading: boolean;
  currentUser: User | null;
  memberships: Membership[];
  activeBusinessId: string | null;
  activeMembership: Membership | null;
  login: () => Promise<void>;
  logout: () => Promise<void>;
  refreshSession: () => Promise<void>;
  setActiveBusinessId: (businessId: string) => void;
}

const ACTIVE_BUSINESS_STORAGE_KEY = 'active_business_id';

const AuthContext = createContext<AuthContextType | undefined>(undefined);

const readStoredActiveBusinessId = () => localStorage.getItem(ACTIVE_BUSINESS_STORAGE_KEY);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const [memberships, setMemberships] = useState<Membership[]>([]);
  const [activeBusinessId, setActiveBusinessIdState] = useState<string | null>(readStoredActiveBusinessId());

  const applySession = (session: CurrentUserResponse) => {
    setCurrentUser(session.user);
    setMemberships(session.memberships);
    setIsAuthenticated(true);

    const storedActiveBusinessId = readStoredActiveBusinessId();
    const allowedBusinessIds = new Set(session.memberships.map((membership) => membership.business_id));
    const nextActiveBusinessId = storedActiveBusinessId && allowedBusinessIds.has(storedActiveBusinessId)
      ? storedActiveBusinessId
      : session.active_business_id || session.memberships[0]?.business_id || null;

    if (nextActiveBusinessId) {
      localStorage.setItem(ACTIVE_BUSINESS_STORAGE_KEY, nextActiveBusinessId);
    } else {
      localStorage.removeItem(ACTIVE_BUSINESS_STORAGE_KEY);
    }

    setActiveBusinessIdState(nextActiveBusinessId);
  };

  const clearSession = () => {
    setIsAuthenticated(false);
    setCurrentUser(null);
    setMemberships([]);
    setActiveBusinessIdState(null);
    localStorage.removeItem(ACTIVE_BUSINESS_STORAGE_KEY);
  };

  const refreshSession = async () => {
    setIsLoading(true);
    try {
      const session = await api.getCurrentUser();
      applySession(session);
    } catch {
      clearSession();
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    refreshSession();
  }, []);

  const login = async () => {
    await refreshSession();
  };

  const logout = async () => {
    try {
      await api.logout();
    } catch {
      // best-effort logout; clear local state regardless
    }
    clearSession();
    setIsLoading(false);
  };

  const setActiveBusinessId = (businessId: string) => {
    setActiveBusinessIdState(businessId);
    localStorage.setItem(ACTIVE_BUSINESS_STORAGE_KEY, businessId);
  };

  const activeMembership = useMemo(() => {
    return memberships.find((membership) => membership.business_id === activeBusinessId) || null;
  }, [activeBusinessId, memberships]);

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated,
        isLoading,
        currentUser,
        memberships,
        activeBusinessId,
        activeMembership,
        login,
        logout,
        refreshSession,
        setActiveBusinessId,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
