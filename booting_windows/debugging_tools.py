import os
import IPython

HEX_RADIX = 16

class Reverser(object):
  def __init__( self,
                live_entry_point,
                image_base,
                image_entry_point ):
    self.live_image_offset = live_entry_point - image_entry_point
    self.file_image_base   = image_base

  def addr2ImageAddr( self, liveAddr ):
    imageAddr = liveAddr - self.live_image_offset + self.file_image_base
    # print( hex( imageAddr ) )
    return imageAddr, hex( imageAddr )

  def imageAddr2Addr( self, imageAddr ):
    liveAddr = imageAddr - self.file_image_base + self.live_image_offset
    # print( hex( liveAddr ) )
    return liveAddr, hex( liveAddr )

  def advanceToAddress( self, imageDumpFile, desiredLiveAddress ):
    disassemblyLine = None
    addr            = None

    desiredAddress, _ = self.addr2ImageAddr( desiredLiveAddress )
    print( "Will search for code starting at %s" % hex( desiredAddress ) )

    while( True ):
      try:
        disassemblyLine       = imageDumpFile.readline().strip()
        addrStr, instruction  = disassemblyLine.split( ":" ) # May throw ValueError
        addr                  = int( addrStr, HEX_RADIX )
        if addr == desiredAddress:
          break

      except ValueError:
        pass #skip lines which are not assembly instructions

    return disassemblyLine, addr

  def generateBreakpoinCmdsOnCalls( self, imageDisasPath, liveStartAddress ):
    with open( imageDisasPath ) as imageDumpFile:
      disassemblyLine, currentAddr = \
          self.advanceToAddress( imageDumpFile, liveStartAddress )

      print( "Starting analysis at: %s" % disassemblyLine )

      breakpointCmds = []

      while "ret" not in disassemblyLine:
        disassemblyLine = imageDumpFile.readline().strip()
        if "call" in disassemblyLine:
          print( "Found call at: %s" % disassemblyLine )

          #Add a breakpoint at call site
          addrStr,_ = disassemblyLine.split( ":" )
          liveAddr,_ = self.imageAddr2Addr( int( addrStr, HEX_RADIX ) )
          breakpointCmds.append( "b *%s" % hex( liveAddr ) )

          # Add a breakpoint immediately after call
          disassemblyLine = imageDumpFile.readline().strip()
          try:
            addrStr = disassemblyLine.split( ":" )[0]
          except:
            print( disassemblyLine )
            raise
          liveAddr,_ = self.imageAddr2Addr( int( addrStr, HEX_RADIX ) )
          breakpointCmds.append( "b *%s" % hex( liveAddr ) )

      print( "Reached end of function at: %s" % disassemblyLine )

      for bp in breakpointCmds:
        print( bp )

      return bp

  def breakpointFromImageAddr( self, imageAddr ):
    addr,_ = self.imageAddr2Addr( imageAddr )
    print( "b *%s" % hex( addr ) )

def main():
  IPython.embed()

if __name__ == '__main__':
  main()

