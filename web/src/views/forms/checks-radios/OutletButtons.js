import React, { useEffect, useState } from 'react'
import { CCard, CCardBody, CCardHeader, CCol, CFormCheck, CFormSwitch, CRow } from '@coreui/react'
import { DocsComponents, DocsExample } from 'src/components'
import { setOutlet, getOutlets, getOutletState, getOutletSysInfo } from '../../../api/outlets'

const ChecksRadios = () => {
  const [outletStates, setOutletStates] = useState({});
  const [outlets, setOutlets] = useState([]);
  const [outletNames, setOutletNames] = useState({});
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setIsLoading(true);
        const discoveredOutlets = await getOutlets();
        console.log('Discovered outlets:', discoveredOutlets);

        if (discoveredOutlets?.result?.ips) {
          setOutlets(discoveredOutlets.result.ips);

          const initialStates = {};
          const initialNames = {};
          for (const id of discoveredOutlets.result.ips) {
            const state = await getOutletState('kasa', id);
            const sysInfo = await getOutletSysInfo('kasa', id);
            initialStates[id] = state.result.state === "True";
            initialNames[id] = sysInfo.result.alias || `Office Outlet (${id})`;
          }
          setOutletStates(initialStates);
          setOutletNames(initialNames);
        } else {
          console.error('Invalid response format:', discoveredOutlets);
          setOutlets([]);
        }
      } catch (error) {
        console.error('Error fetching outlets:', error);
        setOutlets([]);
      } finally {
        setIsLoading(false);
      }
    }
    fetchData();
  }, []);

  if (isLoading) {
    return <div>Loading...</div>;
  }

  const ToggleOutlet = (id) => async () => {
    const action = outletStates[id] ? 'off' : 'on';
    await setOutlet('kasa', id, action);
    setOutletStates(prev => ({
      ...prev,
      [id]: !prev[id]
    }));
  }

  return (
    <CRow>
      <CCol xs={12}>
        <CCard className="mb-4">
          <CCardHeader>
            <strong>Outlets</strong>
          </CCardHeader>
          <CCardBody>
            {outlets.map((id) => (
              <CFormSwitch
                key={id}
                label={outletNames[id] || `Outlet (${id})`}
                id={`formSwitch-${id}`}
                checked={!!outletStates[id]}
                onChange={ToggleOutlet(id)}
              />
            ))}
          </CCardBody>
        </CCard>
      </CCol>
    </CRow>
  )
}

export default ChecksRadios
